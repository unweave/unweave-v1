package config

import "github.com/unweave/unweave-v2/session/model"

type DBConfig struct {
	Host     string `json:"host" env:"UNWEAVE_DB_HOST"`
	Port     int    `json:"port" env:"UNWEAVE_DB_PORT"`
	Name     string `json:"name" env:"UNWEAVE_DB_NAME"`
	User     string `json:"user" env:"UNWEAVE_DB_USER"`
	Password string `json:"password" env:"UNWEAVE_DB_PASSWORD"`
}

type LambdaLabsConfig struct {
	APIKey string `json:"apiKey" env:"UNWEAVE_LAMBDALABS_API_KEY"`
	SSHKey string `json:"sshKey" env:"UNWEAVE_LAMBDALABS_SSH_KEY"`
}

type SessionConfig struct {
	Runtime    model.RuntimeProvider `json:"runtime"`
	LambdaLabs LambdaLabsConfig      `json:"lambdalabs"`
}

type Config struct {
	APIPort string   `json:"port" env:"UNWEAVE_API_PORT"`
	DB      DBConfig `json:"db"`
}
