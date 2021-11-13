package domain

type Zone struct {
	Id      int    `json:"id" db:"id"`
	EventId int    `json:"eventId" db:"eventid"`
	Taken   string `json:"taken" db:"taken"`
	Price   int    `json:"price" db:"price"`
}
