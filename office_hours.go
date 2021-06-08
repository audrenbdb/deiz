package deiz

import "time"

type OfficeHours struct {
	ID          int         `json:"id"`
	StartMn     int         `json:"startMn"`
	EndMn       int         `json:"endMn"`
	WeekDay     int         `json:"weekDay"`
	Address     Address     `json:"address"`
	MeetingMode MeetingMode `json:"meetingMode"`
}

type MeetingMode uint8

const (
	RemoteAvailability MeetingMode = iota
	InOfficeAvailability
	AtExternalAddress
)

func (h *OfficeHours) IsValid() bool {
	return h.StartMn < h.EndMn && h.WeekDay >= 0 && h.WeekDay <= 6
}

func (h *OfficeHours) IsInvalid() bool {
	return !h.IsValid()
}

func (h *OfficeHours) IsWithinDate(d time.Time) bool {
	return int(d.Weekday()) == h.WeekDay
}
