package deiz

type OfficeHours struct {
	ID      int     `json:"id" validator:"required"`
	StartMn int     `json:"startMn" validator:"required"`
	EndMn   int     `json:"endMn" validator:"required"`
	WeekDay int     `json:"weekDay" validator:"required"`
	Address Address `json:"address" validator:"required"`
}
