package deiz

import (
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

func (b *Booking) PatientNotSet() bool {
	return b.Patient.ID == 0
}

func (b *Booking) ClinicianNotSet() bool {
	return b.Clinician.ID == 0
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

func (b *Booking) SetBlocked() {
	b.Patient.ID = 0
	b.Address.ID = 0
	b.Motive.ID = 0
	b.Note = ""
	b.Blocked = true
}
