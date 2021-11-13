package repository

import (
	"context"
	"errors"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable        = "users"
	adminsTable       = "admins"
	zonesTable        = "zones"
	eventsTable       = "events"
	sessionsTable     = "sessions"
	verificationTable = "verification"
	ordersTable       = "orders"
)

var (
	errUserAlreadyExists       = errors.New("user with such email already exists")
	errUserNotFound            = errors.New("user doesn't exists")
	errVerificationCodeInvalid = errors.New("verification code is invalid")
)

type Admins interface {
	GetByCredentials(ctx context.Context, email, passwordHash string) (string, error)
	SetSession(ctx context.Context, id string, session domain.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (string, error)
	IsDuplicate(email string) bool

	IsAdmin(ctx context.Context, id string) (bool, error)
	AddNewAdmin(ctx context.Context, email, passwordHash string) error
	CreateEvent(ctx context.Context, event domain.Event) error
	AddZones(ctx context.Context, eventId int, zones []domain.Zone) error
	DeleteEvent(ctx context.Context, id int) error
	DeleteUser(ctx context.Context, id string) error
	ChangeEvent(ctx context.Context, id int, event domain.Event) error
}

type Users interface {
	CreateUser(ctx context.Context, user domain.User) (string, error)
	GetUserEmail(ctx context.Context, id string) (string, error)
	GetByCredentials(ctx context.Context, email, passwordHash string) (string, error)
	SetSession(ctx context.Context, id string, session domain.Session) error
	SetVerCode(ctx context.Context, id, code string) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (string, error)
	Verify(ctx context.Context, id, hash string) error
	DeleteIfNotVer(id string)
}

type Events interface {
	GetEventById(ctx context.Context, id int) (domain.Event, error)
	GetAllEvents(ctx context.Context) ([]domain.Event, error)
	GetFirstEvent(ctx context.Context) (domain.Event, error)
}

type Zones interface {
	GetZonesByEventId(ctx context.Context, id int) ([]domain.Zone, error)
	TakeZonesById(ctx context.Context, idEvent int, idZones []int, userId string) ([]domain.Zone, error)
	GetZonesByUserId(ctx context.Context, userId string) ([]domain.Zone, error)
}

type Orders interface {
}

type Repositories struct {
	Orders
	Admins
	Users
	Events
	Zones
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		NewOrdersRepo(db),
		NewAdminsRepo(db),
		NewUsersRepo(db),
		NewEventsRepo(db),
		NewZonesRepo(db),
	}
}
