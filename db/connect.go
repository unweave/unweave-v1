package db

import (
	"database/sql"
	"fmt"
)

// Q is a querier for the global database connection
//
// This needs to be initialized when the db connection is first established. It is then
// safe to use across go routines.
var Q Querier

type Config struct {
	Host     string `json:"host" env:"UNWEAVE_DB_HOST"`
	Port     int    `json:"port" env:"UNWEAVE_DB_PORT"`
	Name     string `json:"name" env:"UNWEAVE_DB_NAME"`
	User     string `json:"user" env:"UNWEAVE_DB_USER"`
	Password string `json:"password" env:"UNWEAVE_DB_PASSWORD"`
}

func Connect(cfg Config) (*sql.DB, error) {
	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	var err error
	conn, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, err
}
