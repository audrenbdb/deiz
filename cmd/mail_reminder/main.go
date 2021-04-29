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
	"time"
)

func main() {
	ctx := context.Background()

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
		//Client:    mail.NewGmailClient(),
		Client: mail.NewPostFixClient(),
		Intl:   intl.NewIntlParser("Fr", paris),
	})
	reminder := booking.SendReminderUsecase{
		Getter: repo,
		Mailer: mail,
	}
	if err := reminder.SendReminders(ctx); err != nil {
		log.Println(err)
	}
}
