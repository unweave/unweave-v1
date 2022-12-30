package types

const (
	RuntimeProviderKey = "RuntimeProvider"
)

// swagger:enum Status
type Status string

const (
	StatusInitializingNode Status = "initializingNode"
	StatusRunning          Status = "running"
	StatusStoppingNode     Status = "stoppingNode"
)

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

type Node struct {
	ID      string `json:"id"`
	KeyPair SSHKey `json:"sshKeyPair"`
	Status  Status `json:"status"`
}

type SSHKey struct {
	Name       *string `json:"name,omitempty"`
	PrivateKey *string `json:"privateKey,omitempty"`
	PublicKey  *string `json:"publicKey,omitempty"`
}
type Session struct {
	ID     string `json:"id"`
	SSHKey SSHKey `json:"sshKey"`
	Status Status `json:"runtimeStatus"`
}
type ExecParams struct {
	Cmd   []string `json:"cmd"`
	Image []string `json:"containerImage"`
}
