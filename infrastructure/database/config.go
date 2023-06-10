package database

import (
	"os"
)

type config struct {
	host     string
	database string
	port     string
	user     string
	password string
}

func newConfigPostgresql() *config {
	return &config{
		host:     os.Getenv("POSTGRESQL_HOST"),
		database: os.Getenv("POSTGRESQL_DATABASE"),
		port:     os.Getenv("POSTGRESQL_PORT"),
		user:     os.Getenv("POSTGRESQL_USER"),
		password: os.Getenv("POSTGRESQL_PASSWORD"),
	}
}
