package deiz

//BookingMotive represents a reason to consult a professional
//A clinician may have multiple motive with different duration and prices
//PublicRole means that this motive can be selected by a patient checking clinician availabilities
type BookingMotive struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	//Duration in mn
	Duration int   `json:"duration"`
	Price    int64 `json:"price"`
	Public   bool  `json:"public"`
}
