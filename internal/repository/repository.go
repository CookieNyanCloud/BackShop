package repository

import (
	"context"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	usersTable        = "users"
	zonesTable        = "zones"
	eventsTable       = "events"
	sessionsTable     = "sessions"
	verificationTable = "verification"
)

type Admins interface {
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	SetSession(ctx context.Context, userId uuid.UUID, session domain.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	CreateEvent(time time.Time, description string, zones []domain.Zone) (int, error)
}

type Users interface {
	CreateUser(ctx context.Context, user domain.User) (string, error)
	IsDuplicate(email string) bool
	GetUserEmail(ctx context.Context, id string) (string, error)
	GetByCredentials(ctx context.Context, email, passwordHash string) (string, error)
	SetSession(ctx context.Context, id string, session domain.Session) error
	SetVerCode(ctx context.Context, id, code string) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (string, error)
	Verify(ctx context.Context, id, hash string) error
	DeleteIfNotVer(id string)
}

type Events interface {
	GetEventById(id int) (domain.Event, error)
	GetAllEvents() ([]domain.Event, error)
}

type Zones interface {
	GetZonesByEventId(id int) ([]domain.Zone, error)
	TakeZonesById(idEvent int, idZones []int, userId string) ([]domain.Zone, error)
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
