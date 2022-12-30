package types

const (
	RuntimeProviderKey = "RuntimeProvider"
)

type SSHKey struct {
	Name       *string `json:"name,omitempty"`
	PrivateKey *string `json:"privateKey,omitempty"`
	PublicKey  *string `json:"publicKey,omitempty"`
}

// swagger:enum Status
type Status string

const (
	StatusInitializingNode Status = "initializingNode"
	StatusRunning          Status = "running"
	StatusStoppingNode     Status = "stoppingNode"
)

type Node struct {
	ID      string `json:"id"`
	KeyPair SSHKey `json:"sshKeyPair"`
	Status  Status `json:"status"`
}

// RuntimeProvider is the platform that the node is spawned on. This is where the user
// code runs
type RuntimeProvider string

func (r RuntimeProvider) String() string {
	return string(r)
}

const (
	LambdaLabsProvider RuntimeProvider = "lambdalabs"
	UnweaveProvider    RuntimeProvider = "unweave"
)
