package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/unweave/unweave/providers/lambdalabs"
	"github.com/unweave/unweave/types"
)

type Initializer interface {
	FromUser(userID uuid.UUID, provider types.RuntimeProvider) (*Runtime, error)
	FromSession(sessionID string, provider types.RuntimeProvider) (*Runtime, error)
}

// ConfigFileInitializer is only used in development or if you're self-hosting Unweave.
type ConfigFileInitializer struct{}

type lambdaLabsRuntimeConfig struct {
	APIKey string `json:"apiKey"`
}

type runtimeConfig struct {
	LambdaLabs lambdaLabsRuntimeConfig `json:"lambdaLabs"`
}

func (i *ConfigFileInitializer) FromUser(userID uuid.UUID, provider types.RuntimeProvider) (*Runtime, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filepath.Join(home, ".unweave/runtime-config.json"))
	if err != nil {
		return nil, err
	}

	var config runtimeConfig
	if err = json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	switch provider {
	case types.LambdaLabsProvider:
		if config.LambdaLabs.APIKey == "" {
			return nil, fmt.Errorf("missing LambdaLabs API key in runtime config file")
		}
		sess, err := lambdalabs.NewSessionProvider(config.LambdaLabs.APIKey)
		if err != nil {
			return nil, err
		}
		return &Runtime{sess}, nil

	case types.UnweaveProvider:
		return nil, fmt.Errorf("unweave provider not supported in config file initializer")
	}
	return nil, fmt.Errorf("unknown runtime provider %q", provider)
}

func (i *ConfigFileInitializer) FromSession(sessionID string, provider types.RuntimeProvider) (*Runtime, error) {
	return nil, nil
}
