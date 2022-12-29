package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/unweave/unweave-v2/pkg/gonfig"
)

type ProjectConfig struct {
	ID         string `json:"id"`
	Token      string `json:"token"`
	SSHKeyPath string `json:"SSHKeyPath" env:"UW_SSH_KEY_PATH"`
}

type UserConfig struct {
	Token string `json:"token"`
}

type Config struct {
	UwEnv    string                   `json:"unweaveEnv" env:"UW_ENV"`
	ApiURL   string                   `json:"apiURL" env:"UW_API_URL"`
	AppURL   string                   `json:"appURL" env:"UW_APP_URL"`
	User     UserConfig               `json:"user"`
	Projects map[string]ProjectConfig `json:"projects"`
}

var Path = ""
var UnweaveConfig = &Config{}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not get user home directory")
	}

	env := os.Getenv("UW_ENV")
	switch env {
	case "staging", "stg":
		Path = filepath.Join(home, ".unweave/stg-config.json")
	case "development", "dev":
		Path = filepath.Join(home, ".unweave/dev-config.json")
	case "production", "prod":
		Path = filepath.Join(home, ".unweave/config.json")
	default:
		// If anything else, assume production
		fmt.Println("Unrecognized environment. Assuming production.")
	}

	// Load saved config - create the empty config if it doesn't exist
	if err = readAndUnmarshal(Path, UnweaveConfig); os.IsNotExist(err) {
		err = marshalAndWrite(Path, UnweaveConfig)
		if err != nil {
			log.Fatal("Failed to create config file: ", err)
		}
	} else if err != nil {
		log.Fatal("Failed to read config file: ", err)
	}

	// Override with environment variables
	gonfig.GetFromEnvVariables(UnweaveConfig)
}
