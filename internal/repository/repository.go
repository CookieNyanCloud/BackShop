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
	CreateUser(ctx context.Context, user domain.User) (uuid.UUID, error)
	IsDuplicate(email, name string) bool
	//GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetUserInfo(ctx context.Context, id uuid.UUID) (domain.User, error)
	SetSession(ctx context.Context, userId uuid.UUID, session domain.Session) error
	SerVerCode(ctx context.Context, code string) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Verify(ctx context.Context, userId uuid.UUID, hash string) error
	DeleteIfNotVer(id uuid.UUID)
}

type Events interface {
	GetEvent() ([]domain.Event, error)
}

type Zones interface {
	GetZonesByEventId(id int) ([]domain.Zone, error)
	TakeZoneById(idEvent, idZone, userId uuid.UUID) ([]domain.Zone, error)
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
