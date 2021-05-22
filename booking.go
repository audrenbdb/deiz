package deiz

import (
	"sort"
	"time"
)

type Booking struct {
	ID int `json:"id"`
	//Description is an optional parameters to attach description details
	//Usually comes with an event that one wishes to describe such as its motive.
	Description string    `json:"description"`
	DeleteID    string    `json:"deleteId"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Clinician   Clinician `json:"clinician"`
	Patient     Patient   `json:"patient"`
	Address     string    `json:"address"`
	Paid        bool      `json:"paid"`
	Confirmed   bool      `json:"confirmed"`
	Note        string    `json:"note"`
	Price       int64     `json:"price"`
	//Title of the booking
	//Can either be :
	//"block", "appointment", "event"
	BookingType BookingType `json:"bookingType"`
	//AvailabilityType to status how is the public available to book
	//Can either be remote / in office / at patient home
	AvailabilityType AvailabilityType `json:"availabilityType"`
}

type BookingType int32

const (
	BlockedBooking BookingType = iota
	AppointmentBooking
	EventBooking
)

func (b *Booking) ToEvent() {
	b.BookingType = EventBooking
	b.Patient.ID = 0
	b.Price = 0
	b.AvailabilityType = AtExternalAddress
}

func (b *Booking) EventValid() bool {
	if b.BookingType != EventBooking {
		return false
	}
	if b.PatientSet() {
		return false
	}
	if b.ClinicianNotSet() {
		return false
	}
	if b.Price != 0 {
		return false
	}
	if b.AvailabilityType != AtExternalAddress {
		return false
	}
	if b.Start.After(b.End) {
		return false
	}
	return true
}

func (b *Booking) Assigned() bool {
	return b.Confirmed || b.PreRegistered()
}

func (b *Booking) PreRegistered() bool {
	return b.ID != 0 && !b.Confirmed
}

func (b *Booking) PatientNotSet() bool {
	return b.Patient.ID == 0
}

func (b *Booking) PatientSet() bool {
	return !b.PatientNotSet()
}

func (b *Booking) ClinicianNotSet() bool {
	return b.Clinician.ID == 0
}

func (b *Booking) ClinicianSet() bool {
	return !b.ClinicianNotSet()
}

func (b *Booking) EndBeforeStart() bool {
	return b.End.Before(b.Start)
}

func (b *Booking) Remote() bool {
	return b.Address == ""
}

func (b *Booking) SetPatient(p Patient) {
	b.Patient = p
}

func (b *Booking) SetBlocked() *Booking {
	b.Patient.ID = 0
	b.Note = ""
	b.BookingType = BlockedBooking
	return b
}

func SortBookingByDate(bookings []Booking) []Booking {
	sort.SliceStable(bookings, func(i, j int) bool {
		return bookings[i].Start.Before(bookings[j].Start)
	})
	return bookings
}
