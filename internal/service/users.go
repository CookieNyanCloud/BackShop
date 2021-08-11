package service

import (
	"context"
	"errors"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/cookienyancloud/back/internal/repository"
	"github.com/cookienyancloud/back/pkg/auth"
	"github.com/cookienyancloud/back/pkg/hash"
	"github.com/cookienyancloud/back/pkg/otp"
	"strconv"
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
		Name:     input.Name,
		Password: passwordHash,
		Email:    input.Email,
		Verification: domain.Verification{
			Code: verificationCode,
		},
	}
	if _, err := s.repo.CreateUser(ctx, user); err != nil {
		if err == repository.ErrUserAlreadyExists {
			return ErrUserAlreadyExists
		}
		return err
	}
	//return s.emailService.SendUserVerificationEmail(VerificationEmailInput{
	//	Email:            user.Email,
	//	Name:             user.Name,
	//	VerificationCode: verificationCode,
	//})
	return nil
}

func (s *UsersService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}
	user, err := s.repo.GetByCredentials(ctx, input.Email, passwordHash)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return Tokens{}, ErrUserNotFound
		}
	}
	IdStr:= strconv.Itoa(user.ID)
	return s.createSession(ctx, IdStr)
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}
	IdStr:= strconv.Itoa(user.ID)
	return s.createSession(ctx, IdStr)
}

func (s *UsersService) createSession(ctx context.Context, userId string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId, s.accessTokenTTL)
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

	err = s.repo.SetSession(ctx, userId, session)

	return res, err
}

func (s *UsersService) GetUserInfo(ctx context.Context, id int) (domain.User, error) {
	user, err := s.repo.GetUserInfo(ctx, id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return domain.User{}, ErrUserNotFound
		}
	}
	return user, nil
}

func (s *UsersService) Verify(ctx context.Context, userId , hash string) error {
	err := s.repo.Verify(ctx, userId, hash)
	if err != nil {
		if errors.Is(err, repository.ErrVerificationCodeInvalid) {
			return err
		}
		return err
	}
	return nil
}