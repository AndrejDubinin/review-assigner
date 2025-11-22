package app

import "time"

type (
	Options struct {
		Host            string
		Port            string
		ReadTimeout     string
		WriteTimeout    string
		IdleTimeout     string
		ShutdownTimeout string
	}
	path struct {
		index string
	}
	web struct {
		port            string
		host            string
		readTimeout     time.Duration
		writeTimeout    time.Duration
		idleTimeout     time.Duration
		shutdownTimeout time.Duration
	}

	config struct {
		web  web
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

	return config{
		web: web{
			port:            opts.Port,
			host:            opts.Host,
			readTimeout:     readTimeout,
			writeTimeout:    writeTimeout,
			idleTimeout:     idleTimeout,
			shutdownTimeout: shutdownTimeout,
		},
		path: path{
			index: "/",
		},
	}, nil
}
