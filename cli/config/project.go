package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func IsProjectLinked() bool {
	_, err := os.Stat(projectConfigPath)
	return err == nil
}

func WriteProjectConfig(projectID uuid.UUID, providers []string) error {
	fmt.Println(providers)
	buf := &bytes.Buffer{}

	quotedProviders := make([]string, len(providers))
	for i, val := range providers {
		quotedProviders[i] = fmt.Sprintf("\"%s\"", val)
	}

	vars := struct {
		ProjectID string
		Providers string
	}{
		ProjectID: projectID.String(),
		Providers: strings.Join(quotedProviders, " ,"),
	}
	if err := configTemplate.Execute(buf, vars); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(projectConfigPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(projectConfigPath, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func WriteEnvConfig() error {
	buf := &bytes.Buffer{}
	if err := envTemplate.Execute(buf, nil); err != nil {
		return err
	}
	dir := filepath.Dir(envConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(envConfigPath, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func WriteGitignore() error {
	dir := filepath.Dir(projectConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(path, []byte(gitignoreEmbed), 0644); err != nil {
		return err
	}
	return nil
}
