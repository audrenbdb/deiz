package deiz

type OfficeHours struct {
	ID      int     `json:"id"`
	StartMn int     `json:"startMn"`
	EndMn   int     `json:"endMn"`
	WeekDay int     `json:"weekDay"`
	Address Address `json:"address"`
}
