package main

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz/gcalendar"
	"github.com/audrenbdb/deiz/gmaps"
	"github.com/audrenbdb/deiz/mail"
	"github.com/audrenbdb/deiz/repo/psql"
	"github.com/audrenbdb/deiz/usecase/booking"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

func main() {
	ctx := context.Background()

	path, err := getPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get path")
		os.Exit(1)
	}
	psqlDB, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to start db pool: %v\n", err)
		os.Exit(1)
	}

	paris, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to load location: %v\n", err)
		os.Exit(1)
	}

	repo := psql.NewRepo(psqlDB, nil)
	gCal := gcalendar.NewService()
	gMaps := gmaps.NewService()
	mail := mail.NewService(parseEmailTemplates(path), mail.NewPostFixClient(), paris)
	//mail := mail.NewService(parseEmailTemplates(path), mail.NewGmailClient(), paris)
	bookingUsecase := booking.NewUsecase(repo, mail, gMaps, gCal, paris)
	bookingUsecase.SendReminders(ctx)

}

func getPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func parseEmailTemplates(path string) *template.Template {
	return template.Must(template.ParseGlob(path + "/../../mail/templates/booking-reminder.html"))
}
