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
	AvailabilityType AvailabilityType  `json:"availabilityType"`
	Recurrence       BookingRecurrence `json:"recurrence"`
}

type BookingType uint8
type BookingRecurrence uint8

const (
	NoRecurrence BookingRecurrence = iota
	DailyRecurrence
	WeeklyRecurrence
	MonthlyRecurrence
)

const (
	BlockedBooking BookingType = iota
	AppointmentBooking
	EventBooking
)

type Notification struct {
	ToPatient   bool
	ToClinician bool
}

func (b *Booking) setAsEvent() {
	b.Patient.ID = 0
	b.Price = 0
	b.AvailabilityType = AtExternalAddress
}

func (b *Booking) IsValid(clinicianID int) bool {
	if b.Start.After(b.End) {
		return false
	}
	if b.Clinician.ID != clinicianID || b.ClinicianNotSet() {
		return false
	}
	switch b.BookingType {
	case BlockedBooking:
		return b.blockedBookingValid()
	case AppointmentBooking:
		return b.appointmentBookingValid()
	case EventBooking:
		return b.eventBookingValid()
	default:
		return true
	}
}

func (b *Booking) IsInvalid(clinicianID int) bool {
	return !b.IsValid(clinicianID)
}

func (b *Booking) eventBookingValid() bool {
	b.setAsEvent()
	if b.PatientSet() {
		return false
	}
	return true
}

func (b *Booking) appointmentBookingValid() bool {
	if b.Confirmed {
		return b.PatientSet()
	}
	return true
}

func (b *Booking) blockedBookingValid() bool {
	b.block()
	return b.PatientNotSet() && len(b.Note) == 0
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

func (b *Booking) block() {
	b.Patient.ID = 0
	b.Note = ""
}

func SortBookingByDate(bookings []Booking) []Booking {
	sort.SliceStable(bookings, func(i, j int) bool {
		return bookings[i].Start.Before(bookings[j].Start)
	})
	return bookings
}
