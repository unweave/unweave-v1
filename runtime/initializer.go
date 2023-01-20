package runtime

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/providers/lambdalabs"
	"github.com/unweave/unweave/tools/gonfig"
)

type Initializer interface {
	FromAccount(ctx context.Context, accountID uuid.UUID, provider types.RuntimeProvider) (*Runtime, error)
	FromToken(ctx context.Context, token string, provider types.RuntimeProvider) (*Runtime, error)
}

// EnvInitializer is only used in development or if you're self-hosting Unweave.
type EnvInitializer struct{}

type envInitializer struct {
	LambdaLabsAPIKey string `env:"LAMBDALABS_API_KEY"`
}

func (i *EnvInitializer) FromAccount(ctx context.Context, userID uuid.UUID, provider types.RuntimeProvider) (*Runtime, error) {
	var config envInitializer
	gonfig.GetFromEnvVariables(&config)

	switch provider {
	case types.LambdaLabsProvider:
		if config.LambdaLabsAPIKey == "" {
			return nil, fmt.Errorf("missing LambdaLabs API key in runtime config file")
		}
		return i.FromToken(ctx, config.LambdaLabsAPIKey, provider)
	case types.UnweaveProvider:
		return nil, fmt.Errorf("unweave provider not supported in the env initializer")
	}
	return nil, fmt.Errorf("unknown runtime provider %q", provider)
}

func (i *EnvInitializer) FromToken(ctx context.Context, token string, provider types.RuntimeProvider) (*Runtime, error) {
	switch provider {
	case types.LambdaLabsProvider:
		sess, err := lambdalabs.NewSessionProvider(token)
		if err != nil {
			return nil, err
		}
		return &Runtime{sess}, nil
	case types.UnweaveProvider:
		return nil, fmt.Errorf("unweave provider not supported in the env initializer")
	}
	return nil, fmt.Errorf("unknown runtime provider %q", provider)
}
