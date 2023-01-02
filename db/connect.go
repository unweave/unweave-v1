package db

import (
	"database/sql"
	"fmt"

	"github.com/unweave/unweave-v2/config"
)

func Connect(cfg config.DBConfig) (*sql.DB, error) {
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
