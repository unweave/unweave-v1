package config

import (
	_ "embed"
	"log"
	"text/template"

	"github.com/pelletier/go-toml/v2"
)

type (
	user struct {
		Token string `toml:"token"`
	}

	sshkey struct {
		Path string `toml:"path"`
		Name string `toml:"name"`
	}

	lambda struct {
		ApiKey string `toml:"api_key"`
	}

	providerSecrets struct {
		LambdaLabs lambda `toml:"lambda_labs"`
	}

	secrets struct {
		Token           string          `toml:"token" env:"UNWEAVE_PROJECT_TOKEN"`
		SshKeys         []sshkey        `toml:"ssh_keys"`
		ProviderSecrets providerSecrets `toml:"provider_secrets"`
	}

	provider struct {
		Default string   `toml:"default"`
		Active  []string `toml:"active"`
	}

	project struct {
		ID       string   `toml:"project_id"`
		Env      secrets  `toml:"env"`
		Provider provider `toml:"providers"`
	}

	config struct {
		UnwEnv  string  `toml:"unweave_env" env:"UNWEAVE_ENV"`
		ApiURL  string  `toml:"api_url" env:"UNWEAVE_API_URL"`
		AppURL  string  `toml:"app_url" env:"UNWEAVE_APP_URL"`
		User    user    `toml:"user"`
		Project project `toml:"project"`
	}
)

var (
	//go:embed templates/config.toml
	configEmbed       string
	configTemplate, _ = template.New("config.toml").Parse(configEmbed)

	//go:embed templates/env.toml
	envEmbed       string
	envTemplate, _ = template.New("env.toml").Parse(envEmbed)

	//go:embed templates/gitignore
	gitignoreEmbed string

	unweaveConfigPath = ""
	projectConfigPath = ".unweave/config.toml"
	envConfigPath     = ".unweave/env.toml"

	Config = &config{
		ApiURL: "https://api.unweave.io",
		AppURL: "https://app.unweave.io",
	}
)

func (c *config) String() string {
	buf, err := toml.Marshal(c)
	if err != nil {
		log.Fatal("Failed to marshal config: ", err)
	}
	return string(buf)
}

func (c *config) Save() error {
	return marshalAndWrite(unweaveConfigPath, c)
}

func (c *project) String() string {
	buf, err := toml.Marshal(c)
	if err != nil {
		log.Fatal("Failed to marshal config: ", err)
	}
	return string(buf)
}

func (c *project) Save() error {
	return marshalAndWrite(projectConfigPath, c)
}
