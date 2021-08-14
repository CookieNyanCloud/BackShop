package domain

import "time"

type User struct {
	ID           int          `json:"-" db:"id"`
	Name         string       `form:"username" json:"username" binding:"required" db:"name"`
	Email        string       `form:"email" json:"email" binding:"required" db:"email"`
	Password     string       `form:"password" json:"password" binding:"required" db:"password_hash"`
	Verification Verification `json:"verification"`
	TakenZones   []Zone       `json:"zone"`
	Session      Session      `json:"session"`
	RegisteredAt time.Time    `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt  time.Time    `json:"lastVisitAt" bson:"lastVisitAt"`
}

type Verification struct {
	Code     string `db:"codes"`
	Verified bool   `json:"verified"`
}
