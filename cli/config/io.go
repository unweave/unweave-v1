package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	} else if err != nil {
		return err
	}
	return nil
}

// readAndUnmarshal reads the config file and unmarshals it into the Config struct
func readAndUnmarshal(path string, config *Config) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, config)
}

// marshalAndWrite marshals a RootConfig struct and writes it to disk. It reloads the
// config variable after writing.
func marshalAndWrite(path string, config *Config) error {
	if err := createDir(filepath.Dir(path)); err != nil {
		return err
	}
	buf, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(path, buf, os.ModePerm); err != nil {
		return err
	}
	return nil
}
