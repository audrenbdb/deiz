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
	"github.com/audrenbdb/deiz/mail/mailtmpl"
	"github.com/audrenbdb/deiz/patient"
	"github.com/audrenbdb/deiz/pdf"
	"github.com/audrenbdb/deiz/repo/psql"
	"github.com/audrenbdb/deiz/stripe"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"path/filepath"
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
	emailTemplates, err := mailtmpl.Embed()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to embed email templates : %v\n", err)
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
		Templates: emailTemplates,
		Client:    mail.NewPostFixClient(),
		//Client: mail.NewGmailClient(),
		Intl: intl,
	})
	err = echo.StartEchoServer(echo.EchoServerDeps{
		ContactService: contact.NewUsecase(repo, mail),
		//CredentialsGetter: echo.FakeCredentialsGetter, //http.FirebaseCredentialsGetter(fbClient),
		CredentialsGetter: echo.FirebaseCredentialsGetter(fbClient),
		AccountUsecases:   newAccountUsecases(repo),
		PatientUsecases:   newPatientUsecases(repo),
		BookingUsecases:   newBookingUsecases(paris, repo, mail),
		BillingUsecases:   newBillingUsecases(repo, mail, pdf),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to start echo http_server: %v\n", err)
		os.Exit(1)
	}
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

func newAccountUsecases(repo *psql.Repo) usecase.AccountUsecases {
	crypt := crypt.NewService()
	motiveUc := &motive.BookingMotiveUsecase{
		MotiveUpdater: repo,
		MotiveDeleter: repo,
		MotiveCreater: repo,
	}
	officeHoursUc := &officehours.Usecase{
		Deleter: repo,
		Creater: repo,
	}
	clinicianUc := &clinician.EditUsecase{
		PhoneUpdater:      repo,
		EmailUpdater:      repo,
		AdeliUpdater:      repo,
		ProfessionUpdater: repo,
	}
	return usecase.AccountUsecases{
		AccountDataGetter: &account.GetDataUsecase{AccountDataGetter: repo},
		LoginAllower: &account.AllowLoginUsecase{
			ClinicianGetter: repo,
			AuthChecker:     repo,
			AuthEnabler:     repo,
		},
		AccountAdder: &account.AddAccountUsecase{AccountCreater: repo},
		ClinicianUsecases: usecase.ClinicianUsecases{
			ProfessionEditer: clinicianUc,
			PhoneEditer:      clinicianUc,
			EmailEditer:      clinicianUc,
			AdeliEditer:      clinicianUc,
		},
		BusinessUsecases: &business.UpdateUsecase{
			BusinessUpdater: repo,
		},
		AccountAddressUsecases: usecase.AccountAddressUsecases{
			OfficeAddressAdder: &address.AddAddressUsecase{AddressCreater: repo},
			AddressDeleter:     &address.DeleteAddressUsecase{AccountGetter: repo, AddressDeleter: repo},
			HomeAddressSetter:  &address.SetHomeUsecase{HomeAddressSetter: repo},
			AddressEditer:      &address.EditAddressUsecase{AddressUpdater: repo, AccountGetter: repo},
		},
		MotiveUsecases: usecase.MotiveUsecases{
			MotiveAdder:   motiveUc,
			MotiveRemover: motiveUc,
			MotiveEditer:  motiveUc,
		},
		OfficeHoursUsecases: usecase.OfficeHoursUsecases{
			OfficeHoursAdder:   officeHoursUc,
			OfficeHoursRemover: officeHoursUc,
		},
		CalendarSettingsUsecases: &settings.CalendarSettingsUsecase{
			SettingsUpdater: repo,
		},
		StripeKeysUsecases: &stripekeys.Usecase{
			Crypter:           crypt,
			StripeKeysUpdater: repo,
		},
	}
}

func newPatientUsecases(repo *psql.Repo) usecase.PatientUsecases {
	uc := &patient.Usecase{
		Searcher:              repo,
		Creater:               repo,
		Updater:               repo,
		GetterByEmail:         repo,
		ClinicianBoundChecker: repo,
		AddressCreater:        repo,
		AddressUpdater:        repo,
		BookingsGetter:        repo,
	}
	return usecase.PatientUsecases{
		Searcher:       uc,
		Adder:          uc,
		Editer:         uc,
		AddressEditer:  uc,
		AddressAdder:   uc,
		BookingsGetter: uc,
	}
}

func newBillingUsecases(repo *psql.Repo, mailer *mail.Mailer, pdf *pdf.Pdf) usecase.BillingUsecases {
	stripe := stripe.NewService()
	crypt := crypt.NewService()
	return usecase.BillingUsecases{
		InvoiceCreater: &billing.CreateInvoiceUsecase{
			Counter:    repo,
			Saver:      repo,
			PdfCreater: pdf,
			Mailer:     mailer,
		},
		InvoiceCanceler: &billing.CancelInvoiceUsecase{
			Counter: repo,
			Saver:   repo,
		},
		InvoiceMailer: &billing.MailInvoiceUsecase{
			InvoiceMailer:             mailer,
			InvoicesSummaryMailer:     mailer,
			PdfInvoicesSummaryCreater: pdf,
			PdfInvoiceCreater:         pdf,
			InvoicesGetter:            repo,
		},
		InvoicesGetter: &billing.GetPeriodInvoicesUsecase{Getter: repo},
		StripeSessionCreater: &billing.CreateStripeSessionUsecase{
			Crypter:              crypt,
			StripeSessionCreater: stripe,
			SecretKeyGetter:      repo,
		},
		UnpaidBookingsGetter: &billing.GetUnpaidBookingsUsecase{Getter: repo},
	}
}

func newBookingUsecases(paris *time.Location, repo *psql.Repo, mailer *mail.Mailer) usecase.BookingUsecases {
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
	return usecase.BookingUsecases{
		Register:       bookingRegister,
		PreRegister:    bookingPreRegister,
		CalendarReader: calendarReader,
		SlotDeleter:    bookingSlotDeleter,
		SlotBlocker:    bookingSlotBlocker,
	}
}
