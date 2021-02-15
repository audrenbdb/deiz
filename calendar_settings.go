package deiz

type CalendarSettings struct {
	ID            int           `json:"id"`
	DefaultMotive BookingMotive `json:"defaultMotive"`
	Timezone      Timezone      `json:"timezone"`
}

type Timezone struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
