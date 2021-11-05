package domain

import "time"

type Event struct {
	Id          int       `json:"-" db:"id"`
	Time        time.Time `json:"time" db:"time"`
	Description string    `json:"description" db:"description"`
	MapFile     string    `json:"mapfile" db:"mapfile"`
	Zones       []Zone
}
