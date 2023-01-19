package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/unweave/unweave/tools/gonfig"
)

type (
	user struct {
		Token string `toml:"token"`
	}

	config struct {
		UnwEnv  string  `toml:"unweave_env" env:"UNWEAVE_ENV"`
		ApiURL  string  `toml:"api_url" env:"UNWEAVE_API_URL"`
		AppURL  string  `toml:"app_url" env:"UNWEAVE_APP_URL"`
		User    user    `toml:"user"`
		Project project `toml:"project"`
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
	return marshalAndWrite(UnweaveConfigPath, c)
}

var (
	UnweaveConfigPath = ""
	ProjectConfigPath = ".unweave/config.toml"
	Config            = &config{
		ApiURL: "https://api.unweave.io",
		AppURL: "https://app.unweave.io",
	}
)

func init() {
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
		UnweaveConfigPath = filepath.Join(home, ".unweave/stg-config.toml")
		apiURL = "https://api.staging-unweave.io"
		appURL = "https://app.staging-unweave.io"
	case "development", "dev":
		UnweaveConfigPath = filepath.Join(home, ".unweave/dev-config.toml")
		apiURL = "http://localhost:4000"
		appURL = "http://localhost:3000"
	case "production", "prod":
		UnweaveConfigPath = filepath.Join(home, ".unweave/config.toml")
	default:
		// If anything else, assume production
		fmt.Println("Unrecognized environment. Assuming production.")
	}

	// Load saved config - create the empty config if it doesn't exist
	if err = readAndUnmarshal(UnweaveConfigPath, Config); os.IsNotExist(err) {
		err = marshalAndWrite(UnweaveConfigPath, Config)
		if err != nil {
			log.Fatal("Failed to create config file: ", err)
		}
	} else if err != nil {
		log.Fatal("Failed to read config file: ", err)
	}

	// Need to set these after reading the config file so that they can be overridden
	Config.ApiURL = apiURL
	Config.AppURL = appURL

	// Override with environment variables
	gonfig.GetFromEnvVariables(Config)
}
