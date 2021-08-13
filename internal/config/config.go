package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"time"
)

const (
	defaultHttpPort               = "8000"
	defaultHttpRWTimeout          = 10 * time.Second
	defaultHttpMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
	defaultLimiterRPS             = 10
	defaultLimiterBurst           = 2
	defaultLimiterTTL             = 10 * time.Minute
	defaultVerificationCodeLength = 8
)

type (
	Config struct {
		Environment string
		Postgres    PostgresConfig
		HTTP        HTTPConfig
		Auth        AuthConfig
		Email       EmailConfig
		Limiter     LimiterConfig
		CacheTTL    time.Duration
		SMTP        SMTPConfig
		Payment     PaymentConfig
	}

	PostgresConfig struct {
		Host     string
		Port     string
		Username string
		DBName   string
		SSLMode  string
		Password string
	}

	AuthConfig struct {
		JWT                    JWTConfig
		PasswordSalt           string
		VerificationCodeLength int
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration
		RefreshTokenTTL time.Duration
		SigningKey      string
	}

	EmailConfig struct {
		Templates EmailTemplates
		Subjects  EmailSubjects
	}

	EmailTemplates struct {
		Verification       string `mapstructure:"verification_email"`
		PurchaseSuccessful string `mapstructure:"purchase_successful"`
	}

	EmailSubjects struct {
		Verification       string `mapstructure:"verification_email"`
		PurchaseSuccessful string `mapstructure:"purchase_successful"`
	}

	HTTPConfig struct {
		Host               string
		Port               string
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		MaxHeaderMegabytes int
	}

	LimiterConfig struct {
		RPS   int
		Burst int
		TTL   time.Duration
	}

	SMTPConfig struct {
		Host string
		Port int
		Pass string
		From string
	}

	PaymentConfig struct {
		Fondy       FondyConfig
		CallbackURL string
		ResponseURL string
	}

	FondyConfig struct {
		MerchantId       string
		MerchantPassword string
	}
)

func Init(configsDir string) (*Config, error) {
	populateDefaults()

	if err := parseEnv(); err != nil {
		return nil, err
	}

	if err := parseConfigFile(configsDir, viper.GetString("env")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	setFromEnv(&cfg)
	println(cfg.Postgres.DBName)
	println(cfg.Postgres.SSLMode)
	println(cfg.Postgres.Host)
	println(cfg.Postgres.Password)
	println(cfg.Postgres.Username)
	println(cfg.Postgres.Port)
	return &cfg, nil
}

func parseEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := parsePostgresEnvVariables(); err != nil {
		return err
	}
	if err := parseJWTFromEnv(); err != nil {
		return err
	}
	if err := parseHostFromEnv(); err != nil {
		return err
	}
	if err := parseSMTPPassFromEnv(); err != nil {
		return err
	}
	if err := parseFondyEnvVariables(); err != nil {
		return err
	}
	if err := parseAppEnvFromEnv(); err != nil {
		return err
	}
	if err := parsePaymentEnvVariables(); err != nil {
		return err
	}
	return parsePasswordFromEnv()
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.SetConfigName(env)
	return viper.MergeInConfig()
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("cache.ttl", &cfg.CacheTTL); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth.verificationCodeLength", &cfg.Auth.VerificationCodeLength); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("limiter", &cfg.Limiter); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("smtp", &cfg.SMTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("email.templates", &cfg.Email.Templates); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("email.subjects", &cfg.Email.Subjects); err != nil {
		return err
	}
	return nil
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.Password = viper.GetString("pass")

	cfg.Auth.PasswordSalt = viper.GetString("salt")
	cfg.Auth.JWT.SigningKey = viper.GetString("signing_key")

	cfg.HTTP.Host = viper.GetString("host")

	cfg.SMTP.Pass = viper.GetString("pass")

	cfg.Environment = viper.GetString("env")

	cfg.Payment.Fondy.MerchantId = viper.GetString("merchant_id")
	cfg.Payment.Fondy.MerchantPassword = viper.GetString("merchant_pass")
	cfg.Payment.CallbackURL = viper.GetString("callback_url")
	cfg.Payment.ResponseURL = viper.GetString("redirect_url")

}

func parsePostgresEnvVariables() error {
	viper.SetEnvPrefix("postgres")
	return viper.BindEnv("pass")
}
func parsePasswordFromEnv() error {
	viper.SetEnvPrefix("password")
	return viper.BindEnv("salt")
}
func parseJWTFromEnv() error {
	viper.SetEnvPrefix("jwt")
	return viper.BindEnv("signing_key")
}
func parseHostFromEnv() error {
	viper.SetEnvPrefix("http")
	return viper.BindEnv("host")
}
func parseSMTPPassFromEnv() error {
	viper.SetEnvPrefix("smtp")

	return viper.BindEnv("password")
}
func parseFondyEnvVariables() error {
	viper.SetEnvPrefix("fondy")

	if err := viper.BindEnv("merchant_id"); err != nil {
		return err
	}

	return viper.BindEnv("merchant_pass")
}
func parsePaymentEnvVariables() error {
	viper.SetEnvPrefix("payment")

	if err := viper.BindEnv("callback_url"); err != nil {
		return err
	}

	return viper.BindEnv("redirect_url")
}
func parseAppEnvFromEnv() error {
	viper.SetEnvPrefix("app")
	return viper.BindEnv("env")
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHttpPort)
	viper.SetDefault("http.max_header_megabytes", defaultHttpMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHttpRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHttpRWTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
	viper.SetDefault("auth.verificationCodeLength", defaultVerificationCodeLength)
	viper.SetDefault("limiter.rps", defaultLimiterRPS)
	viper.SetDefault("limiter.burst", defaultLimiterBurst)
	viper.SetDefault("limiter.ttl", defaultLimiterTTL)
}
