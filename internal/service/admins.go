package service

import (
	"context"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/cookienyancloud/back/internal/repository"
	"github.com/cookienyancloud/back/pkg/auth"
	"github.com/cookienyancloud/back/pkg/hash"
	"os"
	"strconv"
	"time"
)

type AdminsService struct {
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	repo       repository.Admins

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAdminsService(
	hasher hash.PasswordHasher,
	tokenManager auth.TokenManager,
	repo repository.Admins,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	) *AdminsService {
	return &AdminsService{
		hasher:          hasher,
		tokenManager:    tokenManager,
		repo:            repo,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AdminsService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}
	admin, err := s.repo.GetByCredentials(ctx, input.Email, passwordHash)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return Tokens{}, ErrUserNotFound
	}
		}
	need:= os.Getenv("ADMINpwd")
	if need!= admin.Email{
		return Tokens{}, NotAdmin
	}
	IdStr:= strconv.Itoa(admin.ID)
	return s.createSession(ctx, IdStr)

}

func (s *AdminsService) createSession(ctx context.Context, adminId string) (Tokens, error) {
	var (
		res Tokens
		err error
	)
	//println("v jwt",userId)
	res.AccessToken, err = s.tokenManager.NewJWT(adminId, s.accessTokenTTL)
	if err != nil {
		return res, err
	}
	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}
	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}
	err = s.repo.SetSession(ctx, adminId, session)
	println("zdes")
	return res, err
}

func (s *AdminsService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}
	IdStr:= strconv.Itoa(user.ID)
	return s.createSession(ctx, IdStr)
}

func (s *AdminsService) CreateEvent(input CreateEventInput) (int, error){
	return s.repo.CreateEvent(input.Time, input.Description, input.Zones)
}