package session

import (
	"github.com/unweave/unweave-v2/config"
	"github.com/unweave/unweave-v2/session/runtime"
)

type Session struct {
	runtime.Runtime
}

func NewSession(cfg config.SessionConfig) Session {
	switch cfg.Runtime {
	case config.LambdaLabs:
		return Session{}
	case config.Unweave:
		return Session{}
	default:
		panic("Unknown runtime provider")
	}
}
