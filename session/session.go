package session

import (
	"github.com/unweave/unweave-v2/session/providers/lambdalabs"
	"github.com/unweave/unweave-v2/session/providers/unweave"
	"github.com/unweave/unweave-v2/types"
)

func NewRuntime(provider types.RuntimeProvider) (Runtime, error) {
	switch provider {
	case types.LambdaLabsProvider:
		return lambdalabs.NewProvider("")

	case types.UnweaveProvider:
		return unweave.NewProvider("")

	default:
		panic("Unknown runtime provider")
	}
}
