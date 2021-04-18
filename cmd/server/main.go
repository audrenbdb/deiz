package main

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz/crypt"
	"github.com/audrenbdb/deiz/gcalendar"
	"github.com/audrenbdb/deiz/gmaps"
	"github.com/audrenbdb/deiz/http"
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

	paris, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to load location: %v\n", err)
		os.Exit(1)
	}

	repo := psql.NewRepo(psqlDB, fbClient)
	pdf := pdf.NewService("oxygen", "oxygen.ttf", filepath.Join(path, "../../assets", "fonts"), paris)
	//mail := mail.NewService(parseEmailTemplates(path), mail.NewPostFixClient(), paris)
	mail := mail.NewService(parseEmailTemplates(path), mail.NewGmailClient(), paris)
	gCal := gcalendar.NewService()
	gMaps := gmaps.NewService()
	stripe := stripe.NewService()
	crypt := crypt.NewService()
	bookingRegister := booking.NewRegisterUsecase(
		paris, repo, repo, repo, repo, repo, mail, gCal, gMaps)
	bookingPreRegister := booking.NewPreRegisterUsecase(repo, repo)

	calendarReader := booking.NewCalendarReaderUsecase(paris, repo, repo)
	bookingSlotDeleter := booking.NewSlotDeleterUsecase(repo, repo, mail)
	bookingSlotBlocker := booking.NewSlotBlockerUsecase(repo)
	err = http.StartEchoServer(
		//http.FirebaseCredentialsGetter(fbClient),
		http.FakeCredentialsGetter,
		account.NewUsecase(repo, crypt),
		patient.NewUsecase(repo),
		billing.NewUsecase(repo, mail, pdf, crypt, stripe),
		contact.NewUsecase(repo, mail),
		bookingRegister,
		bookingPreRegister,
		bookingSlotBlocker,
		bookingSlotDeleter,
		calendarReader,
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
