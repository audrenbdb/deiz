package deiz

type OfficeHours struct {
	ID      int     `json:"id"`
	StartMn int     `json:"startMn"`
	EndMn   int     `json:"endMn"`
	WeekDay int     `json:"weekDay"`
	Address Address `json:"address"`
}

func (h *OfficeHours) IsValid() bool {
	return h.StartMn < h.EndMn && h.WeekDay >= 0 && h.WeekDay <= 6
}

func (h *OfficeHours) IsInvalid() bool {
	return !h.IsValid()
}
