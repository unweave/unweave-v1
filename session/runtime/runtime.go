package runtime

type SSHConnection struct {
	Host     string
	Port     string
	User     string
	Password string
}

type SSHKeyPair struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

// swagger:enum Status
type Status string

const (
	StatusInitializingNode Status = "initializingNode"
	StatusRunning          Status = "running"
	StatusStoppingNode     Status = "stoppingNode"
)

type Node struct {
	ID      string     `json:"id"`
	KeyPair SSHKeyPair `json:"sshKeyPair"`
	Status  Status     `json:"status"`
}

type Runtime interface {
	InitNode() (node Node, err error)
	TerminateNode() error
}
