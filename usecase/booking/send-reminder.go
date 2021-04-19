package booking

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	bookingsInTimeRangeGetter interface {
		GetBookingsInTimeRange(ctx context.Context, start, end time.Time) ([]deiz.Booking, error)
	}
	reminderMailer interface {
		MailBookingReminder(b *deiz.Booking) error
	}
)

type reminder struct {
	getter bookingsInTimeRangeGetter
	mailer reminderMailer
}

type ReminderDeps struct {
	Getter bookingsInTimeRangeGetter
	Mailer reminderMailer
}

func NewReminderUsecase(deps ReminderDeps) *reminder {
	return &reminder{
		getter: deps.Getter,
		mailer: deps.Mailer,
	}
}

func (r *reminder) SendReminders(ctx context.Context) error {
	rangeToFetch := getReminderRange()
	bookings, err := r.getter.GetBookingsInTimeRange(ctx, rangeToFetch.start, rangeToFetch.end)
	if err != nil {
		return fmt.Errorf("unable to get bookings in time range: %s", err)
	}
	for _, b := range bookings {
		err := r.mailer.MailBookingReminder(&b)
		if err != nil {
			return err
		}
	}
	return nil
}

//Get a time range to scan for upcoming bookings in 48h
func getReminderRange() timeRange {
	anchor := time.Now().AddDate(0, 0, 2).UTC()
	start := time.Date(anchor.Year(), anchor.Month(), anchor.Day(), anchor.Hour(), 0, 0, 0, time.UTC)
	return timeRange{
		start: start,
		end:   start.Add(time.Hour * time.Duration(1)).UTC(),
	}
}

type timeRange struct {
	start time.Time
	end   time.Time
}
