package service

import (
	"context"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/cookienyancloud/back/internal/repository"
	"github.com/cookienyancloud/back/pkg/auth"
	"github.com/cookienyancloud/back/pkg/hash"
	"github.com/cookienyancloud/back/pkg/otp"
	"time"
)

type UsersService struct {
	repo                   repository.Users
	hasher                 hash.PasswordHasher
	tokenManager           auth.TokenManager
	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
	otpGenerator           otp.Generator
	emailService           Emails
	verificationCodeLength int
}

func NewUsersService(
	repo repository.Users,
	hasher hash.PasswordHasher,
	tokenManager auth.TokenManager,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	otpGenerator otp.Generator,
	verificationCodeLength int,
	emailService Emails,
) *UsersService {
	return &UsersService{
		repo,
		hasher,
		tokenManager,
		accessTTL,
		refreshTTL,
		otpGenerator,
		emailService,
		verificationCodeLength,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input UserSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)
	user := domain.User{
		Password: passwordHash,
		Email:    input.Email,
	}
	var id string
	if id, err = s.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	if err = s.repo.SetVerCode(ctx, id, verificationCode); err != nil {
		//todo:delete user
		return err
	}

	if err = s.emailService.SendUserVerificationEmail(VerificationEmailInput{
		Email:            user.Email,
		VerificationCode: verificationCode,
	}); err != nil {
		//todo:delete user and code
		return err
	}
	return nil
}

func (s *UsersService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
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

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	id, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}
	return s.createSession(ctx, id)
}

func (s *UsersService) createSession(ctx context.Context, id string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(id, s.accessTokenTTL)
	if err != nil {
		return res, err
	}
	//todo:better refresh
	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, id, session)

	return res, err
}

func (s *UsersService) GetUserEmail(ctx context.Context, id string) (string, error) {
	email, err := s.repo.GetUserEmail(ctx, id)
	if err != nil {
		return "", err
	}
	return email, nil
}

func (s *UsersService) Verify(ctx context.Context, id, hash string) error {
	err := s.repo.Verify(ctx, id, hash)
	if err != nil {
		return err
	}
	return nil
}
