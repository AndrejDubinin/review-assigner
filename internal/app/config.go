// Package app defines configuration structs, parses durations from string options,
// and provides the main App struct for server setup and startup.
//
// The package includes Options for HTTP server configuration (host, port, timeouts),
// internal web and path structs, NewConfig to parse and validate durations,
// and App with methods to initialize the server, register handlers, and start listening.
package app

import (
	"fmt"
	"strconv"
	"time"
)

type (
	Options struct {
		Host            string
		Port            string
		ReadTimeout     string
		WriteTimeout    string
		IdleTimeout     string
		ShutdownTimeout string
		DbName          string
		DbUser          string
		DbPassword      string
		DbHost          string
		DbPort          string
		DbMaxConns      string
		DbMinConns      string
		DbMaxConnLife   string
		DbConnMaxIdle   string
	}
	path struct {
		index   string
		teamAdd string
	}
	web struct {
		port            string
		host            string
		readTimeout     time.Duration
		writeTimeout    time.Duration
		idleTimeout     time.Duration
		shutdownTimeout time.Duration
	}
	db struct {
		dsn         string
		maxConns    int32
		minConns    int32
		maxConnLife time.Duration
		connMaxIdle time.Duration
	}

	config struct {
		web  web
		db   db
		path path
	}
)

func NewConfig(opts Options) (config, error) {
	readTimeout, err := time.ParseDuration(opts.ReadTimeout)
	if err != nil {
		return config{}, err
	}
	writeTimeout, err := time.ParseDuration(opts.WriteTimeout)
	if err != nil {
		return config{}, err
	}
	idleTimeout, err := time.ParseDuration(opts.IdleTimeout)
	if err != nil {
		return config{}, err
	}
	shutdownTimeout, err := time.ParseDuration(opts.ShutdownTimeout)
	if err != nil {
		return config{}, err
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", opts.DbUser, opts.DbPassword,
		opts.DbHost, opts.DbPort, opts.DbName)

	dbMaxConns, err := strconv.ParseInt(opts.DbMaxConns, 10, 32)
	if err != nil {
		return config{}, err
	}
	dbMinConns, err := strconv.ParseInt(opts.DbMinConns, 10, 32)
	if err != nil {
		return config{}, err
	}
	dbMaxConnLife, err := time.ParseDuration(opts.DbMaxConnLife)
	if err != nil {
		return config{}, err
	}
	dbConnMaxIdle, err := time.ParseDuration(opts.DbConnMaxIdle)
	if err != nil {
		return config{}, err
	}

	return config{
		web: web{
			port:            opts.Port,
			host:            opts.Host,
			readTimeout:     readTimeout,
			writeTimeout:    writeTimeout,
			idleTimeout:     idleTimeout,
			shutdownTimeout: shutdownTimeout,
		},
		db: db{
			dsn:         dsn,
			maxConns:    int32(dbMaxConns),
			minConns:    int32(dbMinConns),
			maxConnLife: dbMaxConnLife,
			connMaxIdle: dbConnMaxIdle,
		},
		path: path{
			index:   "/",
			teamAdd: "POST /team/add",
		},
	}, nil
}
