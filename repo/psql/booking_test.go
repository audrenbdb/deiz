package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestInsertBookingIntegration(t *testing.T) {
	ctx := context.Background()
	pool, _ := pgxpool.Connect(ctx, os.Getenv("TEST_DATABASE_URL"))

	start := time.Now()
	end := start.Add(time.Hour * 2)
	b := &deiz.Booking{
		Clinician: deiz.Clinician{ID: 1},
		Patient:   deiz.Patient{ID: 8},
		Start:     start,
		End:       end,
	}
	err := insertBooking(ctx, pool, b)
	assert.NoError(t, err)
	assert.Positive(t, b.ID)
}

func TestDeleteBookingIntegration(t *testing.T) {
	ctx := context.Background()
	pool, _ := pgxpool.Connect(ctx, os.Getenv("TEST_DATABASE_URL"))

	b := &deiz.Booking{
		Clinician: deiz.Clinician{ID: 1},
		Patient:   deiz.Patient{ID: 8},
		Start:     time.Now(),
		End:       time.Now().Add(1),
	}
	insertBooking(ctx, pool, b)

	err := deleteBooking(ctx, pool, b.ID, 1)
	assert.NoError(t, err)
}

func TestGetBookingSlotsInTimeRangeIntegration(t *testing.T) {
	ctx := context.Background()
	pool, _ := pgxpool.Connect(ctx, os.Getenv("TEST_DATABASE_URL"))

	start := time.Now()
	end := time.Now().Add(1)

	b := &deiz.Booking{
		Clinician: deiz.Clinician{ID: 1},
		Patient:   deiz.Patient{ID: 8},
		Start:     start,
		End:       end,
	}
	insertBooking(ctx, pool, b)

	repo := repo{conn: pool}
	s, err := repo.GetBookingsInTimeRange(ctx, start, end, 1)
	assert.NoError(t, err)
	assert.NotNil(t, s)
}
