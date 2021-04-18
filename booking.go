package deiz

import (
	"sort"
	"time"
)

type Booking struct {
	ID        int           `json:"id"`
	DeleteID  string        `json:"deleteId"`
	Start     time.Time     `json:"start"`
	End       time.Time     `json:"end"`
	Motive    BookingMotive `json:"motive"`
	Clinician Clinician     `json:"clinician"`
	Patient   Patient       `json:"patient"`
	Address   Address       `json:"address"`
	Remote    bool          `json:"remote"`
	Paid      bool          `json:"paid"`
	Blocked   bool          `json:"blocked"`
	Confirmed bool          `json:"confirmed"`
	Note      string        `json:"note"`
}

func (b *Booking) Assigned() bool {
	return b.Confirmed || b.PreRegistered()
}

func (b *Booking) PreRegistered() bool {
	return !b.Blocked && b.ID != 0 && !b.Confirmed
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

func (b *Booking) AddressNotSet() bool {
	return b.Address.ID == 0
}

func (b *Booking) EndBeforeStart() bool {
	return b.End.Before(b.Start)
}

func (b *Booking) RemoteStatusMatchAddress() bool {
	if b.Remote {
		return b.AddressNotSet()
	}
	return true
}

func (b *Booking) SetPatient(p Patient) {
	b.Patient = p
}

func (b *Booking) SetBlocked() *Booking {
	b.Patient.ID = 0
	b.Address.ID = 0
	b.Motive.ID = 0
	b.Note = ""
	b.Blocked = true
	return b
}

func SortBookingByDate(bookings []Booking) []Booking {
	sort.SliceStable(bookings, func(i, j int) bool {
		return bookings[i].Start.Before(bookings[j].Start)
	})
	return bookings
}
