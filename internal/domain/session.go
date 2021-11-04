package domain

import "time"

type Session struct {
	Id           string    `json:"userId" db:"id"`
	RefreshToken string    `json:"refreshToken" db:"refreshtoken"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expiresat"`
	LastVisitAt  time.Time `json:"lastVisitAt" db:"lastvisitat"`
}
