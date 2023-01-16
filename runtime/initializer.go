package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/unweave/unweave/providers/lambdalabs"
	"github.com/unweave/unweave/types"
)

type Initializer interface {
	Initialize(ctx context.Context, provider types.RuntimeProvider) (*Runtime, error)
}

// ConfigFileInitializer is only used in development or if you're self-hosting Unweave.
type ConfigFileInitializer struct {
	Path string
}

type runtimeConfig struct {
	LambdaLabs struct {
		APIKey string `json:"apiKey"`
	} `json:"lambdaLabs"`
}

func (i *ConfigFileInitializer) Initialize(ctx context.Context, provider types.RuntimeProvider) (*Runtime, error) {
	f, err := os.Open(i.Path)
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
