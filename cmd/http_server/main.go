package main

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz/account"
	"github.com/audrenbdb/deiz/account/address"
	"github.com/audrenbdb/deiz/account/business"
	"github.com/audrenbdb/deiz/account/clinician"
	"github.com/audrenbdb/deiz/account/motive"
	"github.com/audrenbdb/deiz/account/officehours"
	"github.com/audrenbdb/deiz/account/settings"
	"github.com/audrenbdb/deiz/account/stripekeys"
	"github.com/audrenbdb/deiz/billing"
	"github.com/audrenbdb/deiz/booking"
	"github.com/audrenbdb/deiz/contact"
	"github.com/audrenbdb/deiz/crypt"
	"github.com/audrenbdb/deiz/http/echo"
	"github.com/audrenbdb/deiz/intl"
	"github.com/audrenbdb/deiz/mail"
	"github.com/audrenbdb/deiz/patient"
	"github.com/audrenbdb/deiz/pdf"
	"github.com/audrenbdb/deiz/repo/psql"
	"github.com/audrenbdb/deiz/stripe"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"path/filepath"
	"text/template"
	"time"

	firebase "firebase.google.com/go"
	firebaseAuth "firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	path, err := getPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get path: %v\n", err)
		os.Exit(1)
	}

	fbClient, err := newFireBaseClient(ctx, path+"/firebase.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to start new fire base client: %v\n", err)
		os.Exit(1)
	}

	psqlDB, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to start database pool : %v\n", err)
		os.Exit(1)
	}
	repo := psql.NewRepo(psqlDB, fbClient)
	paris, _ := time.LoadLocation("Europe/Paris")
	intl := intl.NewIntlParser("Fr", paris)
	pdf := pdf.NewService(pdf.ServiceDeps{
		FontFamily: "oxygen",
		FontFile:   "oxygen.ttf",
		FontDir:    filepath.Join(path, "../../assets", "fonts"),
		Intl:       intl,
	})
	mail := mail.NewService(mail.Deps{
		Templates: parseEmailTemplates(path),
		Client:    mail.NewPostFixClient(), //mail.NewGmailClient(),
		Intl:      intl,
	})
	bookingUsecases := newBookingUsecases(paris, repo, mail)
	billingUsecases := newBillingUsecases(repo, mail, pdf)
	accountUsecases := newAccountUsecases(repo)
	err = echo.StartEchoServer(echo.EchoServerDeps{
		CredentialsGetter: echo.FakeCredentialsGetter, //http.FirebaseCredentialsGetter(fbClient),
		AccountUsecases: echo.AccountUsecases{
			AccountAdder:      accountUsecases.accountAdder,
			LoginAllower:      accountUsecases.loginAllower,
			AccountDataGetter: accountUsecases.accountDataGetter,

			AccountAddressUsecases: echo.AccountAddressUsecases{
				OfficeAddressAdder: accountUsecases.addressUsecases.adder,
				AddressDeleter:     accountUsecases.addressUsecases.deleter,
				HomeAddressSetter:  accountUsecases.addressUsecases.homeSetter,
				AddressEditer:      accountUsecases.addressUsecases.editer,
			},
			BusinessUsecases: accountUsecases.businessUsecases.updater,
			ClinicianUsecases: echo.ClinicianUsecases{
				PhoneEditer:      accountUsecases.clinicianUsecases.editer,
				EmailEditer:      accountUsecases.clinicianUsecases.editer,
				ProfessionEditer: accountUsecases.clinicianUsecases.editer,
				AdeliEditer:      accountUsecases.clinicianUsecases.editer,
			},
			MotiveUsecases: echo.MotiveUsecases{
				MotiveAdder:   accountUsecases.motiveUsecases,
				MotiveEditer:  accountUsecases.motiveUsecases,
				MotiveRemover: accountUsecases.motiveUsecases,
			},
			OfficeHoursUsecases: echo.OfficeHoursUsecases{
				OfficeHoursAdder:   accountUsecases.officeHoursUsecases,
				OfficeHoursRemover: accountUsecases.officeHoursUsecases,
			},
			StripeKeysUsecases:       accountUsecases.stripeKeysUsecases,
			CalendarSettingsUsecases: accountUsecases.calendarSettingsUsecases,
		},
		PatientUsecases: &patient.Usecase{
			Searcher:              repo,
			Creater:               repo,
			Updater:               repo,
			GetterByEmail:         repo,
			ClinicianBoundChecker: repo,
			AddressCreater:        repo,
			AddressUpdater:        repo,
			BookingsGetter:        repo,
		},
		ContactService: contact.NewUsecase(repo, mail),
		BookingUsecases: echo.BookingUsecases{
			Register:       bookingUsecases.register,
			PreRegister:    bookingUsecases.preRegister,
			CalendarReader: bookingUsecases.calendarReader,
			SlotDeleter:    bookingUsecases.slotDeleter,
			SlotBlocker:    bookingUsecases.slotBlocker,
		},
		BillingUsecases: echo.BillingUsecases{
			InvoiceCreater:       billingUsecases.invoiceCreater,
			InvoiceMailer:        billingUsecases.invoiceMailer,
			InvoiceCanceler:      billingUsecases.invoiceCanceler,
			InvoicesGetter:       billingUsecases.periodInvoicesGetter,
			StripeSessionCreater: billingUsecases.stripeSessionCreater,
			UnpaidBookingsGetter: billingUsecases.unpaidBookingsGetter,
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to start echo http_server: %v\n", err)
		os.Exit(1)
	}
}

func parseEmailTemplates(path string) *template.Template {
	return template.Must(template.ParseGlob(path + "/../../mail/templates/*.html"))
}

func newFireBaseClient(ctx context.Context, path string) (*firebaseAuth.Client, error) {
	opt := option.WithCredentialsFile(path)
	fbApp, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	return fbApp.Auth(ctx)
}

func getPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

type accountUsecases struct {
	accountDataGetter *account.GetDataUsecase
	loginAllower      *account.AllowLoginUsecase
	accountAdder      *account.AddAccountUsecase

	clinicianUsecases        clinicianUsecases
	businessUsecases         businessUsecases
	addressUsecases          addressUsecases
	motiveUsecases           *motive.BookingMotiveUsecase
	officeHoursUsecases      *officehours.Usecase
	calendarSettingsUsecases *settings.CalendarSettingsUsecase
	stripeKeysUsecases       *stripekeys.Usecase
}

func newAccountUsecases(repo *psql.Repo) accountUsecases {
	return accountUsecases{
		accountDataGetter: &account.GetDataUsecase{AccountDataGetter: repo},
		loginAllower: &account.AllowLoginUsecase{
			ClinicianGetter: repo,
			AuthChecker:     repo,
			AuthEnabler:     repo,
		},
		accountAdder:      &account.AddAccountUsecase{AccountCreater: repo},
		clinicianUsecases: newClinicianUsecases(repo),
		businessUsecases:  newBusinessUsecases(repo),
		addressUsecases:   newAddressUsecases(repo),
		motiveUsecases: &motive.BookingMotiveUsecase{
			MotiveCreater: repo,
			MotiveDeleter: repo,
			MotiveUpdater: repo,
		},
		officeHoursUsecases: &officehours.Usecase{
			Creater: repo,
			Deleter: repo,
		},
		calendarSettingsUsecases: &settings.CalendarSettingsUsecase{
			SettingsUpdater: repo,
		},
		stripeKeysUsecases: &stripekeys.Usecase{
			StripeKeysUpdater: repo,
		},
	}
}

type clinicianUsecases struct {
	editer *clinician.EditUsecase
}

func newClinicianUsecases(repo *psql.Repo) clinicianUsecases {
	return clinicianUsecases{
		editer: &clinician.EditUsecase{
			EmailUpdater:      repo,
			PhoneUpdater:      repo,
			AdeliUpdater:      repo,
			ProfessionUpdater: repo,
		},
	}
}

type businessUsecases struct {
	updater *business.UpdateUsecase
}

func newBusinessUsecases(repo *psql.Repo) businessUsecases {
	return businessUsecases{
		updater: &business.UpdateUsecase{BusinessUpdater: repo},
	}
}

type addressUsecases struct {
	adder      *address.AddAddressUsecase
	deleter    *address.DeleteAddressUsecase
	editer     *address.EditAddressUsecase
	homeSetter *address.SetHomeUsecase
}

func newAddressUsecases(repo *psql.Repo) addressUsecases {
	return addressUsecases{
		adder:   &address.AddAddressUsecase{AddressCreater: repo},
		deleter: &address.DeleteAddressUsecase{AccountGetter: repo, AddressDeleter: repo},
		editer: &address.EditAddressUsecase{
			AccountGetter:  repo,
			AddressUpdater: repo,
		},
		homeSetter: &address.SetHomeUsecase{
			HomeAddressSetter: repo,
		},
	}
}

type billingUsecases struct {
	invoiceCreater       *billing.CreateInvoiceUsecase
	invoiceCanceler      *billing.CancelInvoiceUsecase
	invoiceMailer        *billing.MailInvoiceUsecase
	periodInvoicesGetter *billing.GetPeriodInvoicesUsecase
	stripeSessionCreater *billing.CreateStripeSessionUsecase
	unpaidBookingsGetter *billing.GetUnpaidBookingsUsecase
}

func newBillingUsecases(repo *psql.Repo, mailer *mail.Mailer, pdf *pdf.Pdf) billingUsecases {
	stripe := stripe.NewService()
	crypt := crypt.NewService()
	return billingUsecases{
		invoiceCreater: &billing.CreateInvoiceUsecase{
			Counter:    repo,
			Saver:      repo,
			PdfCreater: pdf,
			Mailer:     mailer,
		},
		invoiceCanceler: &billing.CancelInvoiceUsecase{
			Counter: repo,
			Saver:   repo,
		},
		invoiceMailer: &billing.MailInvoiceUsecase{
			InvoiceMailer:             mailer,
			InvoicesSummaryMailer:     mailer,
			PdfInvoicesSummaryCreater: pdf,
			PdfInvoiceCreater:         pdf,
			InvoicesGetter:            repo,
		},
		periodInvoicesGetter: &billing.GetPeriodInvoicesUsecase{Getter: repo},
		stripeSessionCreater: &billing.CreateStripeSessionUsecase{
			Crypter:              crypt,
			StripeSessionCreater: stripe,
			SecretKeyGetter:      repo,
		},
		unpaidBookingsGetter: &billing.GetUnpaidBookingsUsecase{Getter: repo},
	}
}

type bookingUsecases struct {
	register       *booking.RegisterUsecase
	preRegister    *booking.PreRegisterUsecase
	calendarReader *booking.ReadCalendarUsecase
	slotDeleter    *booking.DeleteSlotUsecase
	slotBlocker    *booking.BlockSlotUsecase
}

func newBookingUsecases(paris *time.Location, repo *psql.Repo, mailer *mail.Mailer) bookingUsecases {
	bookingRegister := &booking.RegisterUsecase{
		Loc:            paris,
		PatientGetter:  repo,
		PatientCreater: repo,
		BookingCreater: repo,
		BookingUpdater: repo,
		BookingGetter:  repo,
		BookingMailer:  mailer,
	}
	bookingPreRegister := &booking.PreRegisterUsecase{
		BookingGetter:  repo,
		BookingCreater: repo,
	}
	calendarReader := &booking.ReadCalendarUsecase{
		Loc:               paris,
		OfficeHoursGetter: repo,
		BookingsGetter:    repo,
	}
	bookingSlotDeleter := &booking.DeleteSlotUsecase{
		BookingGetter:  repo,
		BookingDeleter: repo,
		CancelMailer:   mailer,
	}
	bookingSlotBlocker := &booking.BlockSlotUsecase{
		Blocker: repo,
	}
	return bookingUsecases{
		register:       bookingRegister,
		preRegister:    bookingPreRegister,
		calendarReader: calendarReader,
		slotDeleter:    bookingSlotDeleter,
		slotBlocker:    bookingSlotBlocker,
	}
}
