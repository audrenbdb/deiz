package motive

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	updater interface {
		UpdateBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error
	}
	deleter interface {
		DeleteBookingMotive(ctx context.Context, bookingMotiveID, clinicianID int) error
	}
	creater interface {
		CreateBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error
	}
)

type BookingMotiveUsecase struct {
	MotiveUpdater updater
	MotiveDeleter deleter
	MotiveCreater creater
}

func (u *BookingMotiveUsecase) EditBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error {
	return u.MotiveUpdater.UpdateBookingMotive(ctx, m, clinicianID)
}

func (u *BookingMotiveUsecase) RemoveBookingMotive(ctx context.Context, mID, clinicianID int) error {
	return u.MotiveDeleter.DeleteBookingMotive(ctx, mID, clinicianID)
}

func (u *BookingMotiveUsecase) AddBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error {
	return u.MotiveCreater.CreateBookingMotive(ctx, m, clinicianID)
}
