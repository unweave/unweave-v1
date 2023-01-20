package config

import (
	_ "embed"
	"text/template"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/unweave/unweave/cli/ui"
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
		NodeTypes []string `toml:"node_types"`
	}

	project struct {
		ID              uuid.UUID           `toml:"project_id"`
		Env             *secrets            `toml:"env"`
		Providers       map[string]provider `toml:"provider"`
		DefaultProvider string              `toml:"default_provider"`
	}

	unweave struct {
		UnwEnv string `toml:"unweave_env" env:"UNWEAVE_ENV"`
		ApiURL string `toml:"api_url" env:"UNWEAVE_API_URL"`
		AppURL string `toml:"app_url" env:"UNWEAVE_APP_URL"`
		User   *user  `toml:"user"`
	}

	config struct {
		Unweave *unweave `toml:"unweave"`
		Project *project `toml:"project"`
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
		Unweave: &unweave{
			ApiURL: "https://api.unweave.io",
			AppURL: "https://app.unweave.io",
			User:   &user{},
		},
		Project: &project{
			Env:       &secrets{},
			Providers: map[string]provider{},
		},
	}
)

func (c *config) String() string {
	buf, err := toml.Marshal(c)
	if err != nil {
		ui.Errorf("Failed to marshal config: ", err)
	}
	return string(buf)
}

func (c *unweave) Save() error {
	return marshalAndWrite(unweaveConfigPath, c)
}

func (c *project) String() string {
	buf, err := toml.Marshal(c)
	if err != nil {
		ui.Errorf("Failed to marshal config: ", err)
	}
	return string(buf)
}

func (c *project) Save() error {
	return marshalAndWrite(projectConfigPath, c)
}
