package config

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
)

type (
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
		DefaultProvider string `toml:"default_provider"`
	}

	project struct {
		ID       string   `toml:"project_id"`
		Env      secrets  `toml:"env"`
		Provider provider `toml:"provider"`
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
)

func (c *project) String() string {
	buf, err := toml.Marshal(c)
	if err != nil {
		log.Fatal("Failed to marshal config: ", err)
	}
	return string(buf)
}

func (c *project) Save() error {
	return marshalAndWrite(ProjectConfigPath, c)
}

func WriteProjectConfig(projectID uuid.UUID) error {
	buf := &bytes.Buffer{}
	vars := struct {
		ProjectID string
	}{
		ProjectID: projectID.String(),
	}
	if err := configTemplate.Execute(buf, vars); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(ProjectConfigPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(ProjectConfigPath, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func WriteEnvConfig() error {
	buf := &bytes.Buffer{}
	if err := envTemplate.Execute(buf, nil); err != nil {
		return err
	}
	dir := filepath.Dir(ProjectConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, "env.toml")
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func WriteGitignore() error {
	dir := filepath.Dir(ProjectConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(path, []byte(gitignoreEmbed), 0644); err != nil {
		return err
	}
	return nil
}
