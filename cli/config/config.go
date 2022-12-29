package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/unweave/unweave-v2/pkg/gonfig"
)

// AuthToken is used to authenticate with the Unweave API. It is loaded from the saved
// config file and can be overridden with runtime flags.
var AuthToken = ""

// ProjectPath is the path to the current project to run commands on. It is loaded from the saved
// config file and can be overridden with runtime flags.
var ProjectPath = ""

// SSHKeyPath is the path to the SSH key to use to connect to a new or existing Session.
var SSHKeyPath = ""

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

func (c *Config) String() string {
	buf, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal config: ", err)
	}
	return string(buf)
}

var Path = ""
var UnweaveConfig = &Config{
	ApiURL: "https://api.unweave.io",
	AppURL: "https://app.unweave.io",
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not get user home directory")
	}

	env := "production"
	if e, ok := os.LookupEnv("UW_ENV"); ok {
		env = e
	}
	switch env {
	case "staging", "stg":
		Path = filepath.Join(home, ".unweave/stg-config.json")
		UnweaveConfig.ApiURL = "https://api.staging-unweave.io"
		UnweaveConfig.AppURL = "https://app.staging-unweave.io"
	case "development", "dev":
		Path = filepath.Join(home, ".unweave/dev-config.json")
		UnweaveConfig.ApiURL = "http://localhost:8080"
		UnweaveConfig.AppURL = "http://localhost:3000"
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
