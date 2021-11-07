package app

import (
	"context"
	"errors"
	"github.com/cookienyancloud/back/internal/config"
	delivery "github.com/cookienyancloud/back/internal/delivery/http"
	"github.com/cookienyancloud/back/internal/repository"
	"github.com/cookienyancloud/back/internal/server"
	"github.com/cookienyancloud/back/internal/service"
	"github.com/cookienyancloud/back/pkg/auth"
	//"github.com/cookienyancloud/back/pkg/cache"
	"github.com/cookienyancloud/back/pkg/database/postgres"
	"github.com/cookienyancloud/back/pkg/email/smtp"
	"github.com/cookienyancloud/back/pkg/hash"
	"github.com/cookienyancloud/back/pkg/logger"
	"github.com/cookienyancloud/back/pkg/otp"
	"github.com/cookienyancloud/back/pkg/payment/fondy"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//todo: idempotent api
//todo: cache


func Run(configPath string, local bool) {

	cfg, err := config.Init(configPath, local)
	if err != nil {
		logger.Error(err)
		return
	}
	dataBaseClient, err := postgres.NewClient(cfg.Postgres)
	if err != nil {
		logger.Error(err)
		return
	}

	//memCache := cache.NewMemoryCache()

	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)


	paymentProvider := fondy.NewFondyClient(cfg.Payment.Fondy.MerchantId, cfg.Payment.Fondy.MerchantPassword)

	//emailProvider := smtp.NewClient(memCache)

	emailSender, err := smtp.NewSMTPSender(
		cfg.SMTP.From,
		cfg.SMTP.Pass,
		cfg.SMTP.Host,
		cfg.SMTP.Port)
	if err != nil {
		logger.Error(err)
		return
	}
	println("S")
	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}

	otpGenerator := otp.NewGOTPGenerator()

	repos := repository.NewRepositories(dataBaseClient)

	services := service.NewServices(service.Deps{
		Repos:                  repos,
		Hasher:                 hasher,
		TokenManager:           tokenManager,
		//EmailProvider:          emailProvider,
		EmailSender:            emailSender,
		EmailConfig:            config.EmailConfig{},
		PaymentProvider:        paymentProvider,
		AccessTokenTTL:         cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:        cfg.Auth.JWT.RefreshTokenTTL,
		PaymentCallbackURL:     cfg.Payment.CallbackURL,
		PaymentResponseURL:     cfg.Payment.ResponseURL,
		CacheTTL:               int64(cfg.CacheTTL.Seconds()),
		OtpGenerator:           otpGenerator,
		VerificationCodeLength: cfg.Auth.VerificationCodeLength,
	})
	handlers := delivery.NewHandler(services,tokenManager)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	if err := dataBaseClient.Close(); err != nil {
		logrus.Errorf("error occurred on db connection close: %s", err.Error())
	}

}
