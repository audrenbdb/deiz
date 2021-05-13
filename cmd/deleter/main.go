package main

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz/booking"
	"github.com/audrenbdb/deiz/repo/psql"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	psqlDB, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to start db pool: %v\n", err)
		os.Exit(1)
	}
	repo := psql.NewRepo(psqlDB, nil)
	uc := booking.BlockSlotUsecase{
		Deleter: repo,
	}
	if err := uc.DeletePastBlockedBookingSlot(ctx); err != nil {
		log.Println(err)
	}
}
