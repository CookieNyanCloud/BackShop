package domain

import "time"

type User struct {
	ID           int          `json:"-" db:"id"`
	Name         string       `form:"username" json:"username" binding:"required"`
	Email        string       `form:"email" json:"email" binding:"required"`
	Password     string       `form:"password" json:"password" binding:"required"`
	Verification Verification `json:"verification" `
	TakenZones   []Zone       `json:"session"`
	Session      Session      `json:"session"`
	RegisteredAt time.Time    `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt  time.Time    `json:"lastVisitAt" bson:"lastVisitAt"`
}

type Verification struct {
	Code     string `db:"codes"`
	Verified bool   `json:"verified"`
}
