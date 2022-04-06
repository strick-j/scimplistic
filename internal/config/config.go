package config

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/strick-j/scimplistic/internal/web"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	ListenAddress              string   `envconfig:"SCIMPLISTIC_LISTEN_ADDR" default:":8080" yaml:"listen"`
	ReadTimeout                Duration `envconfig:"SCIMPLISTIC_READ_TIMEOUT" yaml:"read_timeout"`
	WriteTimeout               Duration `envconfig:"SCIMPLISTIC_WRITE_TIMEOUT" yaml:"write_timeout"`
	TokenBucketTTL             Duration `envconfig:"SCIMPLISTIC_TOKEN_BUCKET_TTL" yaml:"token_bucket_ttl"`
	UserRequestsPerSecondLimit float64  `envconfig:"SCIMPLISTIC_USER_RPS_LIMIT" yaml:"user_rps_limit"`
}

func (cfg ServerConfig) ListenParams() web.ListenParams {
	return web.ListenParams{
		Address:            cfg.ListenAddress,
		ReadTimeout:        cfg.ReadTimeout.Duration,
		WriteTimeout:       cfg.WriteTimeout.Duration,
		LimitExpirationTTL: cfg.TokenBucketTTL.Duration,
		ClientRPSQuota:     cfg.UserRequestsPerSecondLimit,
	}
}

type Database struct {
	Address             string `envconfig:"SCIMPLISTIC_DB_ADDRESS" default:"postgres://localhost:5432/SCIMPLISTIC" yaml:"address"`
	MigrationsDirectory string `envconfig:"SCIMPLISTIC_MIGRATIONS_DIR" default:"deployments/db/migrations" yaml:"migrations_dir"`
	VersionTable        string `envconfig:"SCIMPLISTIC_VERSION_TABLE" default:"schema_migrations" yaml:"version_table"`
	SchemaVersion       uint   `envconfig:"SCIMPLISTIC_SCHEMA_VERSION" yaml:"schema_version"`
	SkipMigration       bool   `envconfig:"SCIMPLISTIC_NO_MIGRATION" yaml:"skip_migration"`
}

func (dbs Database) PoolConfig() (*pgxpool.Config, error) {
	cfg, err := pgxpool.ParseConfig(dbs.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid database DSN string: %w", err)
	}

	return cfg, nil
}

type Redis struct {
	DB       int    `envconfig:"SCIMPLISTIC_REDIS_DB" yaml:"db"`
	Address  string `envconfig:"SCIMPLISTIC_REDIS_ADDRESS" yaml:"address" default:"localhost:6379"`
	Username string `envconfig:"SCIMPLISTIC_REDIS_USER" yaml:"username"`
	Password string `envconfig:"SCIMPLISTIC_REDIS_PASSWORD" yaml:"password"`
}

// RedisOptions returns redis connection options
func (r Redis) RedisOptions() *redis.Options {
	return &redis.Options{
		Network:  "tcp",
		Addr:     r.Address,
		DB:       r.DB,
		Username: r.Username,
		Password: r.Password,
	}
}

type Config struct {
	Production bool         `envconfig:"SCIMPLISTIC_PRODUCTION" default:"false" yaml:"production"`
	Server     ServerConfig `yaml:"server"`
	DB         Database     `yaml:"db"`
	Redis      Redis        `yaml:"redis"`
}

func FromFile(cfgPath string) (*Config, error) {
	cfg, err := FromEnv()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	defer f.Close()
	if err = yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %q: %w", cfgPath, err)
	}
	return cfg, nil
}

func FromEnv() (*Config, error) {
	cfg := &Config{}

	// envconfig doesn't work correctly with nested structs and
	// sets invalid env name for nested fields.
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config from environment variables: %w", err)
	}
	return cfg, err
}
