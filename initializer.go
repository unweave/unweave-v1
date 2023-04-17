// This is an example initializer that uses env variables to configure the runtime. You
// can implement custom initializers to support additional providers by implementing the
// runtime.Initializer interface.
package main

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/builder"
	"github.com/unweave/unweave/providers/lambdalabs"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/tools/gonfig"
	"github.com/unweave/unweave/vault"
)

// EnvInitializer is only used in development or if you're self-hosting Unweave.
type EnvInitializer struct{}

type providerConfig struct {
	LambdaLabsAPIKey string `env:"LAMBDALABS_API_KEY"`
}

type builderConfig struct {
	RegistryURI string `env:"UNWEAVE_CONTAINER_REGISTRY_URI"`
}

func (i *EnvInitializer) InitializeRuntime(ctx context.Context, userID string, provider types.Provider) (*runtime.Runtime, error) {
	var cfg providerConfig
	gonfig.GetFromEnvVariables(&cfg)

	switch provider {
	case types.LambdaLabsProvider:
		if cfg.LambdaLabsAPIKey == "" {
			return nil, fmt.Errorf("missing LambdaLabs API key in runtime config file")
		}
		node, err := lambdalabs.NewNodeRuntime(cfg.LambdaLabsAPIKey)
		if err != nil {
			return nil, err
		}
		sess, err := lambdalabs.NewSessionRuntime(cfg.LambdaLabsAPIKey)
		if err != nil {
			return nil, err
		}

		return &runtime.Runtime{Node: node, Session: sess}, nil

	default:
		return nil, fmt.Errorf("%q provider not supported in the env initializer", provider)
	}
}

func (i *EnvInitializer) InitializeBuilder(ctx context.Context, userID string, builderType string) (builder.Builder, error) {
	var cfg builderConfig
	gonfig.GetFromEnvVariables(&cfg)

	if builderType != "docker" {
		return nil, fmt.Errorf("%q builder not supported in the env initializer", builderType)
	}
	logger := &builder.FsLogger{}
	return builder.NewBuilder(logger, cfg.RegistryURI), nil
}

func (i *EnvInitializer) InitializeVault(ctx context.Context) (vault.Vault, error) {
	return vault.NewMemVault(), nil
}
