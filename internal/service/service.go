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
	"time"
)

type UserSignUpInput struct {
	Name     string
	Email    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type VerificationEmailInput struct {
	Email            string
	Name             string
	VerificationCode string
}

type PurchaseSuccessfulEmailInput struct {
	Email   string
	Name    string
	EventId string
	ZoneId  string
}

type CreateEventInput struct {
	Time        time.Time
	Description string
	//MapFile     string        `json:"mapfile" db:"mapfile"`
	Zones []domain.Zone
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	GetUserInfo(ctx context.Context, id int) (domain.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	Verify(ctx context.Context, userId, hash string) error
}

type Admins interface {
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	CreareEvent(ctx context.Context, input CreateEventInput) (int, error)
}

type Events interface {
	GetEvent() ([]domain.Event, error)

}

type Zones interface {
	GetZonesByEventId(id int) ([]domain.Zone, error)
	TakeZoneById(idEvent, idZone, userId int) ([]domain.Zone, error)
}

type Emails interface {
	SendUserVerificationEmail(VerificationEmailInput) error
}

type Services struct {
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
	return &Services{
		adminsService,
		usersService,
		eventsService,
		zonesService,
	}
}
