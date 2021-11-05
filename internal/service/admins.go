package service

import (
	"context"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/cookienyancloud/back/internal/repository"
	"github.com/cookienyancloud/back/pkg/auth"
	"github.com/cookienyancloud/back/pkg/hash"
	"time"
)

type AdminsService struct {
	hasher          hash.PasswordHasher
	tokenManager    auth.TokenManager
	repo            repository.Admins
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
	id, err := s.repo.GetByCredentials(ctx, input.Email, passwordHash)
	if err != nil {
		return Tokens{}, err
	}
	return s.createSession(ctx, id)

}

func (s *AdminsService) createSession(ctx context.Context, adminId string) (Tokens, error) {
	var (
		res Tokens
		err error
	)
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
	return res, err
}

func (s *AdminsService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	id, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}
	return s.createSession(ctx, id)
}

func (s *AdminsService) CreateEvent(ctx context.Context, input CreateEventInput) error {
	return s.repo.CreateEvent(ctx, domain.Event{
		Time:        input.Time,
		Description: input.Description,
		MapFile:     input.MapFile,
		Zones:       input.Zones,
	})
}

func (s *AdminsService)	IsAdmin(ctx context.Context, id string) (bool, error) {
	return s.repo.IsAdmin(ctx,id)
}