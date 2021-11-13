package service

import (
	"context"
	"github.com/cookienyancloud/back/internal/config"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/cookienyancloud/back/internal/repository"
	"github.com/cookienyancloud/back/pkg/auth"
	"github.com/cookienyancloud/back/pkg/cache"
	"github.com/cookienyancloud/back/pkg/email"
	"github.com/cookienyancloud/back/pkg/hash"
	"github.com/cookienyancloud/back/pkg/otp"
	"github.com/cookienyancloud/back/pkg/payment"
	"github.com/google/uuid"
	"time"
)

type UserSignUpInput struct {
	Email    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type CreateEventInput struct {
	Time        time.Time     `json:"time" db:"time"`
	Description string        `json:"description" db:"description"`
	MapFile     string        `json:"mapfile" db:"mapfile"`
	Zones       []domain.Zone `json:"zones" db:"zones"`
}

type VerificationEmailInput struct {
	Email            string
	VerificationCode string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	GetUserEmail(ctx context.Context, id string) (string, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	Verify(ctx context.Context, id, hash string) error
}

type Admins interface {
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	CreateEvent(ctx context.Context, input CreateEventInput) error
	createSession(ctx context.Context, adminId string) (Tokens, error)
	IsAdmin(ctx context.Context, id string) (bool, error)
}

type Events interface {
	GetEvents(ctx context.Context) ([]domain.Event, error)
	GetEventById(ctx context.Context, id int) (domain.Event, error)
	GetFirstEvent(ctx context.Context) (domain.Event, error)
}

type Zones interface {
	GetZonesByEventId(ctx context.Context, id int) ([]domain.Zone, error)
	TakeZonesById(ctx context.Context, idEvent int, idZones []int, userId string) ([]domain.Zone, error)
	GetZonesByUserId(ctx context.Context, userId string) ([]domain.Zone, error)
}

type Emails interface {
	SendUserVerificationEmail(VerificationEmailInput) error
}

type Payments interface {
	GeneratePaymentLink(ctx context.Context, orderId) (string, error)
	ProcessTransaction(ctx context.Context, callback interface{}) error
}

type Orders interface {
	Create(ctx context.Context,orderId uuid.UUID, userId string, eventId int, zonesId []int) error
	AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) (domain.Order, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID, query domain.GetOrdersQuery) ([]domain.Order, int64, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Order, error)
	SetStatus(ctx context.Context, id primitive.ObjectID, status string) error
}

type Services struct {
	Orders
	Payments
	Admins
	Users
	Events
	Zones
}

type Deps struct {
	Repos                  *repository.Repositories
	Cache                  cache.Cache
	Hasher                 hash.PasswordHasher
	TokenManager           auth.TokenManager
	EmailProvider          email.Provider
	EmailSender            email.Sender
	EmailConfig            config.EmailConfig
	PaymentProvider        payment.Provider
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	PaymentCallbackURL     string
	PaymentResponseURL     string
	CacheTTL               int64
	OtpGenerator           otp.Generator
	VerificationCodeLength int
	FrontEndURL            string
}

func NewServices(deps Deps) *Services {
	emailService := NewEmailsService(
		deps.EmailProvider,
		deps.EmailSender,
		deps.EmailConfig,
		deps.FrontEndURL,
	)
	adminsService := NewAdminsService(
		deps.Hasher,
		deps.TokenManager,
		deps.Repos.Admins,
		deps.AccessTokenTTL,
		deps.RefreshTokenTTL,
	)
	usersService := NewUsersService(
		deps.Repos.Users,
		deps.Hasher,
		deps.TokenManager,
		deps.AccessTokenTTL,
		deps.RefreshTokenTTL,
		deps.OtpGenerator,
		deps.VerificationCodeLength,
		emailService,
	)
	zonesService := NewZonesService(
		deps.Repos.Zones,
	)
	eventsService := NewEventsService(
		deps.Repos.Events,
		deps.Repos.Zones,
	)
	OrdersService := NewOrdersService(
		deps.Repos.Orders,
	)
	return &Services{
		OrdersService,
		adminsService,
		usersService,
		eventsService,
		zonesService,
	}
}
