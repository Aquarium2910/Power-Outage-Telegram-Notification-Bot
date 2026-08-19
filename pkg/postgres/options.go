package postgres

import "time"

type Option func(*Postgres)

func MaxPoolSize(size int) Option {
	return func(pg *Postgres) {
		pg.maxPoolSize = size
	}
}

func ConnAttempts(attempts int) Option {
	return func(pg *Postgres) {
		pg.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(pg *Postgres) {
		pg.connTimeout = timeout
	}
}
