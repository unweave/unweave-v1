package session

import (
	"github.com/unweave/unweave-v2/session/model"
	"github.com/unweave/unweave-v2/session/providers/lambdalabs"
	"github.com/unweave/unweave-v2/session/providers/unweave"
)

func NewRuntime(provider model.RuntimeProvider) Runtime {
	switch provider {
	case model.LambdaLabsProvider:
		return lambdalabs.NewProvider("")

	case model.UnweaveProvider:
		return unweave.NewProvider("")

	default:
		panic("Unknown runtime provider")
	}
}
