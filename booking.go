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
