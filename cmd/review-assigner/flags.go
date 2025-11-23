package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AndrejDubinin/review-assigner/internal/app"
)

const (
	defaultPort            = "8080"
	defaultHost            = "0.0.0.0"
	defaultReadTimeout     = "5s"
	defaultWriteTimeout    = "10s"
	defaultIdleTimeout     = "120m"
	defaultShutdownTimeout = "20s"

	defaultDbName        = "review_assigner_db"
	defaultDbUser        = "review_assigner_user"
	defaultDbPassword    = "review_assigner_pass"
	defaultDbHost        = "postgres"
	defaultDbPort        = "5432"
	defaultDbMaxConns    = "25"
	defaultDbMinConns    = "5"
	defaultDbMaxConnLife = "1h"
	defaultDbConnMaxIdle = "30m"
)

var opts = app.Options{}

func initOpts() {
	flag.StringVar(&opts.Host, "host", getEnv("SERVER_HOST", defaultHost),
		fmt.Sprintf("server's host, default: %q", defaultHost))
	flag.StringVar(&opts.Port, "port", getEnv("SERVER_PORT", defaultPort),
		fmt.Sprintf("server's port, default: %q", defaultPort))
	flag.StringVar(&opts.ReadTimeout, "read-timeout", getEnv("SERVER_READ_TIMEOUT", defaultReadTimeout),
		fmt.Sprintf("server's read timeout, default: %q", defaultReadTimeout))
	flag.StringVar(&opts.WriteTimeout, "write-timeout", getEnv("SERVER_WRITE_TIMEOUT", defaultWriteTimeout),
		fmt.Sprintf("server's write timeout, default: %q", defaultWriteTimeout))
	flag.StringVar(&opts.IdleTimeout, "idle-timeout", getEnv("SERVER_IDLE_TIMEOUT", defaultIdleTimeout),
		fmt.Sprintf("server's idle timeout, default: %q", defaultIdleTimeout))
	flag.StringVar(&opts.ShutdownTimeout, "shutdown-timeout", getEnv("SERVER_SHUTDOWN_TIMEOUT", defaultShutdownTimeout),
		fmt.Sprintf("server's shutdown timeout, default: %q", defaultShutdownTimeout))

	flag.StringVar(&opts.DbName, "db-name", getEnv("POSTGRES_DB", defaultDbName),
		fmt.Sprintf("server's database name, default: %q", defaultDbName))
	flag.StringVar(&opts.DbUser, "db-user", getEnv("POSTGRES_USER", defaultDbUser),
		fmt.Sprintf("server's database user, default: %q", defaultDbUser))
	flag.StringVar(&opts.DbPassword, "db-password", getEnv("POSTGRES_PASSWORD", defaultDbPassword),
		fmt.Sprintf("server's database password, default: %q", defaultDbPassword))
	flag.StringVar(&opts.DbHost, "db-host", getEnv("POSTGRES_HOST", defaultDbHost),
		fmt.Sprintf("server's database host, default: %q", defaultDbHost))
	flag.StringVar(&opts.DbPort, "db-port", getEnv("POSTGRES_PORT", defaultDbPort),
		fmt.Sprintf("server's database port, default: %q", defaultDbPort))
	flag.StringVar(&opts.DbMaxConns, "db-max-conns", getEnv("POSTGRES_MAX_CONNS", defaultDbMaxConns),
		fmt.Sprintf("server's database max connections, default: %q", defaultDbMaxConns))
	flag.StringVar(&opts.DbMinConns, "db-min-conns", getEnv("POSTGRES_MIN_CONNS", defaultDbMinConns),
		fmt.Sprintf("server's database min connections, default: %q", defaultDbMinConns))
	flag.StringVar(&opts.DbMaxConnLife, "db-max-conn-life", getEnv("POSTGRES_MAX_CONN_LIFETIME", defaultDbMaxConnLife),
		fmt.Sprintf("server's database max connection life, default: %q", defaultDbMaxConnLife))
	flag.StringVar(&opts.DbConnMaxIdle, "db-conn-max-idle", getEnv("POSTGRES_MAX_CONN_IDLE_TIME", defaultDbConnMaxIdle),
		fmt.Sprintf("server's database max connection idle, default: %q", defaultDbConnMaxIdle))

	flag.Parse()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
