package session

import "github.com/unweave/unweave-v2/config"

type Session struct {
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
