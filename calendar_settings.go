package deiz

type CalendarSettings struct {
	ID            int           `json:"id"`
	DefaultMotive BookingMotive `json:"defaultMotive"`
	Timezone      Timezone      `json:"timezone"`
	RemoteAllowed bool          `json:"remoteAllowed"`
}

type Timezone struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func (s *CalendarSettings) IsValid() bool {
	return s.ID != 0 && s.Timezone.ID != 0
}

func (s *CalendarSettings) IsInvalid() bool {
	return !s.IsValid()
}
