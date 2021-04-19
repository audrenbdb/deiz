package main

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz/crypt"
	"github.com/audrenbdb/deiz/http"
	"github.com/audrenbdb/deiz/intl"
	"github.com/audrenbdb/deiz/mail"
	"github.com/audrenbdb/deiz/pdf"
	"github.com/audrenbdb/deiz/repo/psql"
	"github.com/audrenbdb/deiz/stripe"
	"github.com/audrenbdb/deiz/usecase/account"
	"github.com/audrenbdb/deiz/usecase/billing"
	"github.com/audrenbdb/deiz/usecase/booking"
	"github.com/audrenbdb/deiz/usecase/contact"
	"github.com/audrenbdb/deiz/usecase/patient"
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
	stripe := stripe.NewService()
	crypt := crypt.NewService()
	bookingUsecases := newBookingUsecases(paris, repo, mail)
	err = http.StartEchoServer(
		http.FirebaseCredentialsGetter(fbClient),
		//http.FakeCredentialsGetter,
		account.NewUsecase(repo, crypt),
		patient.NewUsecase(repo),
		billing.NewUsecase(repo, mail, pdf, crypt, stripe),
		contact.NewUsecase(repo, mail),
		http.BookingUsecases{
			Register:       bookingUsecases.register,
			PreRegister:    bookingUsecases.preRegister,
			CalendarReader: bookingUsecases.calendarReader,
			SlotDeleter:    bookingUsecases.slotDeleter,
			SlotBlocker:    bookingUsecases.slotBlocker,
		},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to start echo server: %v\n", err)
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

type bookingUsecases struct {
	register       *booking.Register
	preRegister    *booking.PreRegister
	calendarReader *booking.ReadCalendar
	slotDeleter    *booking.SlotDeleter
	slotBlocker    *booking.SlotBlocker
}

func newBookingUsecases(paris *time.Location, repo *psql.Repo, mailer *mail.Mailer) bookingUsecases {
	bookingRegister := booking.NewRegisterUsecase(booking.RegisterDeps{
		Loc:            paris,
		PatientGetter:  repo,
		PatientCreater: repo,
		BookingCreater: repo,
		BookingUpdater: repo,
		BookingGetter:  repo,
		BookingMailer:  mailer,
	})
	bookingPreRegister := booking.NewPreRegisterUsecase(booking.PreRegisterDeps{
		BookingGetter:  repo,
		BookingCreater: repo,
	})
	calendarReader := booking.NewCalendarReaderUsecase(booking.CalendarReaderDeps{
		Loc:               paris,
		OfficeHoursGetter: repo,
		BookingsGetter:    repo,
	})
	bookingSlotDeleter := booking.NewSlotDeleterUsecase(booking.SlotDeleterDeps{
		BookingGetter:  repo,
		BookingDeleter: repo,
		CancelMailer:   mailer,
	})
	bookingSlotBlocker := booking.NewSlotBlockerUsecase(booking.SlotBlockerDeps{
		Blocker: repo,
	})
	return bookingUsecases{
		register:       bookingRegister,
		preRegister:    bookingPreRegister,
		calendarReader: calendarReader,
		slotDeleter:    bookingSlotDeleter,
		slotBlocker:    bookingSlotBlocker,
	}
}
