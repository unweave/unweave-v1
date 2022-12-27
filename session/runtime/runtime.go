package runtime

type SSHConnection struct {
	Host     string
	Port     string
	User     string
	Password string
}

// swagger:enum Status
type Status string

const (
	StatusInitializingNode Status = "initializingNode"
	StatusRunning          Status = "running"
	StatusStoppingNode     Status = "stoppingNode"
)

type Runtime interface {
	InitNode() (SSHConnection, error)
	StopNode() error
}
