package runtime

type SSHConnection struct {
	Host     string
	Port     string
	User     string
	Password string
}

type Runtime interface {
	InitNode() (SSHConnection, error)
	StopNode() error
}
