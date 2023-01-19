package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/unweave/unweave/cli/ui"
	"github.com/unweave/unweave/tools/gonfig"
)

// getActiveProjectPath returns the active project directory by recursively going up the
// directory tree until it finds a directory that's contains the .unweave/config.toml file
func getActiveProjectPath() (string, error) {
	var activeProjectDir string
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var walk func(path string)
	walk = func(path string) {
		dotUnw := filepath.Join(path, ".unweave")
		if _, err = os.Stat(dotUnw); err == nil {
			if _, err = os.Stat(filepath.Join(dotUnw, "config.toml")); err == nil {
				activeProjectDir = path
				return
			}
		}

		parent := filepath.Dir(path)
		if parent == "." || parent == "/" {
			return
		}
		walk(parent)
	}
	walk(pwd)

	if activeProjectDir == "" {
		return "", fmt.Errorf("no active project found")
	}
	return activeProjectDir, nil
}

func init() {
	// ----- ProjectConfig -----
	// Try loading config from the current directory
	// if not found, try all parent directories
	// if still not found, init empty struct
	// override with environment variables

	projectConfig := project{}
	envConfig := secrets{}
	projectDir, err := getActiveProjectPath()
	if err == nil {
		projectConfigPath = filepath.Join(projectDir, projectConfigPath)
		envConfigPath = filepath.Join(projectDir, envConfigPath)

		if err = readAndUnmarshal(projectConfigPath, &projectConfig); err != nil {
			ui.Infof("Failed to read project config at path %q", projectConfigPath)
		}
		if err = readAndUnmarshal(envConfigPath, &envConfig); err != nil {
			ui.Infof("Failed to read environment config at path %q", envConfigPath)
		}
	}
	projectConfig.Env = envConfig

	// ----- Unweave Config -----
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not get user home directory")
	}

	env := "production"
	if e, ok := os.LookupEnv("UNWEAVE_ENV"); ok {
		env = e
	}
	apiURL := Config.ApiURL
	appURL := Config.AppURL

	switch env {
	case "staging", "stg":
		unweaveConfigPath = filepath.Join(home, ".unweave/stg-config.toml")
		apiURL = "https://api.staging-unweave.io"
		appURL = "https://app.staging-unweave.io"
	case "development", "dev":
		unweaveConfigPath = filepath.Join(home, ".unweave/dev-config.toml")
		apiURL = "http://localhost:4000"
		appURL = "http://localhost:3000"
	case "production", "prod":
		unweaveConfigPath = filepath.Join(home, ".unweave/config.toml")
	default:
		// If anything else, assume production
		fmt.Println("Unrecognized environment. Assuming production.")
	}

	// Load saved config - create the empty config if it doesn't exist
	if err = readAndUnmarshal(unweaveConfigPath, Config); os.IsNotExist(err) {
		err = marshalAndWrite(unweaveConfigPath, Config)
		if err != nil {
			log.Fatal("Failed to create config file: ", err)
		}
	} else if err != nil {
		log.Fatal("Failed to read config file: ", err)
	}

	// Need to set these after reading the config file so that they can be overridden
	Config.ApiURL = apiURL
	Config.AppURL = appURL
	Config.Project = projectConfig

	// Override with environment variables
	gonfig.GetFromEnvVariables(Config)
}
