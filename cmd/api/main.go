package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/xtommas/food-backend/internal/data"
	"github.com/xtommas/food-backend/internal/jsonlog"
)

var (
	buildTime string
	version   string
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	cors struct {
		trustedOrigins []string
	}
	jwt struct {
		secret string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	var cfg config

	// server
	cfg.port = getEnvInt("PORT", 4000, logger)
	cfg.env = getEnv("ENV", "development")

	// database
	cfg.db.dsn = requireEnv("DB_DSN", logger)
	cfg.db.maxOpenConns = getEnvInt("DB_MAX_OPEN_CONNS", 25, logger)
	cfg.db.maxIdleConns = getEnvInt("DB_MAX_IDLE_CONNS", 25, logger)
	cfg.db.maxIdleTime = getEnv("DB_MAX_IDLE_TIME", "15m")

	// rate limiter
	cfg.limiter.rps = getEnvFloat("LIMITER_RPS", 2, logger)
	cfg.limiter.burst = getEnvInt("LIMITER_BURST", 4, logger)
	cfg.limiter.enabled = getEnvBool("LIMITER_ENABLED", true, logger)

	// CORS
	if origins := os.Getenv("CORS_TRUSTED_ORIGINS"); origins != "" {
		cfg.cors.trustedOrigins = strings.Fields(origins)
	}

	// JWT
	cfg.jwt.secret = requireEnv("JWT_SECRET", logger)

	// version
	if os.Getenv("VERSION") == "true" {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func requireEnv(key string, logger *jsonlog.Logger) string {
	val := os.Getenv(key)
	if val == "" {
		logger.PrintFatal(fmt.Errorf("required environment variable %q is not set", key), nil)
	}
	return val
}

func getEnvInt(key string, defaultVal int, logger *jsonlog.Logger) int {
	if val := os.Getenv(key); val != "" {
		n, err := strconv.Atoi(val)
		if err != nil {
			logger.PrintFatal(fmt.Errorf("invalid value for %q: %w", key, err), nil)
		}
		return n
	}
	return defaultVal
}

func getEnvFloat(key string, defaultVal float64, logger *jsonlog.Logger) float64 {
	if val := os.Getenv(key); val != "" {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			logger.PrintFatal(fmt.Errorf("invalid value for %q: %w", key, err), nil)
		}
		return f
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool, logger *jsonlog.Logger) bool {
	if val := os.Getenv(key); val != "" {
		b, err := strconv.ParseBool(val)
		if err != nil {
			logger.PrintFatal(fmt.Errorf("invalid value for %q: %w", key, err), nil)
		}
		return b
	}
	return defaultVal
}
