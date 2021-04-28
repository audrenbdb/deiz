package main

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz/booking"
	"github.com/audrenbdb/deiz/intl"
	"github.com/audrenbdb/deiz/mail"
	"github.com/audrenbdb/deiz/mail/mailtmpl"
	"github.com/audrenbdb/deiz/repo/psql"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	ctx := context.Background()

	/*path*/
	_, err := getPath()
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
	mailTemplates, err := mailtmpl.Embed()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse email templates")
		os.Exit(1)
	}
	repo := psql.NewRepo(psqlDB, nil)
	mail := mail.NewService(mail.Deps{
		Templates: mailTemplates,
		Client:    mail.NewGmailClient(),
		//Client:    mail.NewPostFixClient(),
		Intl: intl.NewIntlParser("Fr", paris),
	})
	//mail := mail.NewService(parseEmailTemplates(path), mail.NewGmailClient(), paris)
	reminder := booking.SendReminderUsecase{
		Getter: repo,
		Mailer: mail,
	}
	if err := reminder.SendReminders(ctx); err != nil {
		log.Println(err)
	}
}

func getPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

/*
func parseEmailTemplates(path string) *template.Template {
	return template.Must(template.ParseGlob(path + "/../../mail/mailtmpl/booking-reminder.html"))
}
*/
