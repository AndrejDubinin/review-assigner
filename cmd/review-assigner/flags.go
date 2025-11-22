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
	flag.Parse()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
