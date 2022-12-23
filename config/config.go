package config

// RuntimeProvider is the platform that the node is spawned on. This is where the user
// code runs
type RuntimeProvider string

const (
	LambdaLabs RuntimeProvider = "lambdalabs"
	Unweave    RuntimeProvider = "unweave"
)

type Config struct {
}

type SessionConfig struct {
	Runtime RuntimeProvider `json:"runtime"`
}
