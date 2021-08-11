package domain

type Zone struct {
	EventId int `json:"eventId" db:"eventid"`
	Id      int `json:"id" db:"id"`
	Taken   int `json:"taken" db:"taken"`
	Price   int `json:"price" db:"price"`
}
