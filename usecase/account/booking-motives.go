package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	BookingMotiveUpdater interface {
		UpdateBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error
	}
	BookingMotiveDeleter interface {
		DeleteBookingMotive(ctx context.Context, bookingMotiveID, clinicianID int) error
	}
	BookingMotiveCreater interface {
		CreateBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error
	}
)

func (u *Usecase) EditBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error {
	return u.BookingMotiveUpdater.UpdateBookingMotive(ctx, m, clinicianID)
}

func (u *Usecase) RemoveBookingMotive(ctx context.Context, mID, clinicianID int) error {
	return u.BookingMotiveDeleter.DeleteBookingMotive(ctx, mID, clinicianID)
}

func (u *Usecase) AddBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error {
	return u.BookingMotiveCreater.CreateBookingMotive(ctx, m, clinicianID)
}
