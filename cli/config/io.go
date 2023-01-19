package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

func createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	} else if err != nil {
		return err
	}
	return nil
}

// readAndUnmarshal reads the config file and unmarshals it into the config struct
func readAndUnmarshal[T any](path string, config *T) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return toml.Unmarshal(buf, config)
}

// marshalAndWrite marshals a RootConfig struct and writes it to disk. It reloads the
// config variable after writing.
func marshalAndWrite[T any](path string, config *T) error {
	if err := createDir(filepath.Dir(path)); err != nil {
		return err
	}
	buf, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	if err = os.WriteFile(path, buf, os.ModePerm); err != nil {
		return err
	}
	return nil
}
