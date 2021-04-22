package booking

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
	"time"
)

func (r *SendReminderUsecase) SendReminders(ctx context.Context) error {
	bookings, err := getBookingsAwaitingRecall(ctx, r.getter)
	if err != nil {
		return err
	}
	for _, b := range bookings {
		err := r.mailer.MailBookingReminder(&b)
		if err != nil {
			return err
		}
	}
	return nil
}

func getBookingsAwaitingRecall(ctx context.Context, getter bookingsInTimeRangeGetter) ([]deiz.Booking, error) {
	rangeToFetch := getReminderRange()
	bookings, err := getter.GetBookingsInTimeRange(ctx, rangeToFetch.start, rangeToFetch.end)
	if err != nil {
		return nil, fmt.Errorf("unable to get bookings in time range: %s", err)
	}
	return bookings, nil
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

type (
	bookingsInTimeRangeGetter interface {
		GetBookingsInTimeRange(ctx context.Context, start, end time.Time) ([]deiz.Booking, error)
	}
	reminderMailer interface {
		MailBookingReminder(b *deiz.Booking) error
	}
)

type SendReminderUsecase struct {
	getter bookingsInTimeRangeGetter
	mailer reminderMailer
}

type ReminderDeps struct {
	Getter bookingsInTimeRangeGetter
	Mailer reminderMailer
}

func NewReminderUsecase(deps ReminderDeps) *SendReminderUsecase {
	return &SendReminderUsecase{
		getter: deps.Getter,
		mailer: deps.Mailer,
	}
}