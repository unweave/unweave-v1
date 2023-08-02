// This is an example initializer that uses env variables to configure the runtime. You
// can implement custom initializers to support additional providers by implementing the
// runtime.Initializer interface.
package main

import (
	"context"
	"fmt"

	"github.com/unweave/unweave-v1/builder"
	"github.com/unweave/unweave-v1/builder/docker"
	"github.com/unweave/unweave-v1/builder/fslogs"
	"github.com/unweave/unweave-v1/tools/gonfig"
	"github.com/unweave/unweave-v1/vault"
)

// EnvInitializer is only used in development or if you're self-hosting Unweave.
type EnvInitializer struct{}

type providerConfig struct {
	LambdaLabsAPIKey string `env:"LAMBDALABS_API_KEY"`
}

type builderConfig struct {
	RegistryURI string `env:"UNWEAVE_CONTAINER_REGISTRY_URI"`
}

func (i *EnvInitializer) InitializeBuilder(ctx context.Context, userID string, builderType string) (builder.Builder, error) {
	var cfg builderConfig
	gonfig.GetFromEnvVariables(&cfg)

	if builderType != "docker" {
		return nil, fmt.Errorf("%q builder not supported in the env initializer", builderType)
	}

	logger := fslogs.NewLogger()

	return docker.NewBuilder(logger, cfg.RegistryURI), nil
}

func (i *EnvInitializer) InitializeVault(ctx context.Context) (vault.Vault, error) {
	return vault.NewMemVault(), nil
}
