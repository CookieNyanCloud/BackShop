package repository

import (
	"context"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	usersTable = "users"
	zonesTable = "zones"
	eventsTable = "events"
	sessionsTable = "sessions"
	codesTable = "codes"
)

type Admins interface {
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	SetSession(ctx context.Context, userId string, session domain.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	CreateEvent(time time.Time, description string, zones []domain.Zone) (int, error)
}

type Users interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	IsDuplicate(email, name string) (bool, error)
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetUserInfo(ctx context.Context, id int) (domain.User, error)
	SetSession(ctx context.Context, userId string, session domain.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Verify(ctx context.Context, userId string , hash string) error
}


type Events interface {
	GetEvent() ([]domain.Event, error)

}

type Zones interface {
	GetZonesByEventId(id int) ([]domain.Zone, error)
	TakeZoneById(idEvent, idZone, userId int) ([]domain.Zone, error)
}


type Repositories struct {
	Admins
	Users
	Events
	Zones
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		NewAdminsRepo(db),
		NewUsersRepo(db),
		NewEventsRepo(db),
		NewZonesRepo(db),
	}
}
