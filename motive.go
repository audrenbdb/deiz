package deiz

import "context"

//BookingMotive represents a reason to consult a professional
//A clinician may have multiple motive with different duration and prices
//Public means that this motive can be selected by a patient checking clinician availabilities
type BookingMotive struct {
	ID   int    `json:"id" validator:"required"`
	Name string `json:"name" validator:"required"`
	//Duration in mn
	Duration int   `json:"duration" validator:"required"`
	Price    int64 `json:"price" validator:"required"`
	Public   bool  `json:"public"`
}

type (
	bookingMotiveAdder interface {
		AddBookingMotive(ctx context.Context, m *BookingMotive, clinicianID int) error
	}
	bookingMotiveRemover interface {
		RemoveBookingMotive(ctx context.Context, m *BookingMotive, clinicianID int) error
	}
)

type (
	//AddBookingMotive add a new clinician booking motive
	AddBookingMotive func(ctx context.Context, m *BookingMotive, clinicianID int) error
	//RemoveBookingMotive removes a clinician booking motive
	RemoveBookingMotive func(ctx context.Context, m *BookingMotive, clinicianID int) error
)

func addBookingMotiveFunc(adder bookingMotiveAdder) AddBookingMotive {
	return func(ctx context.Context, m *BookingMotive, clinicianID int) error {
		return adder.AddBookingMotive(ctx, m, clinicianID)
	}
}

func removeBookingMotiveFunc(remover bookingMotiveRemover) RemoveBookingMotive {
	return func(ctx context.Context, m *BookingMotive, clinicianID int) error {
		return remover.RemoveBookingMotive(ctx, m, clinicianID)
	}
}
