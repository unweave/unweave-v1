package types

import (
	"time"

	"github.com/google/uuid"
)

const (
	RuntimeProviderKey = "RuntimeProvider"
)

type SessionStatus string

const (
	StatusInitializing SessionStatus = "initializing"
	StatusActive       SessionStatus = "active"
	StatusTerminated   SessionStatus = "terminated"
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

func (r RuntimeProvider) DisplayName() string {
	switch r {
	case LambdaLabsProvider:
		return "LambdaLabs"
	case UnweaveProvider:
		return "Unweave"
	default:
		return "Unknown"
	}
}

type NodeSpecs struct {
	VCPUs int `json:"vCPUs"`
	// Memory is the RAM in GB
	Memory int `json:"memory"`
	// GPUMemory is the GPU RAM in GB
	GPUMemory *int `json:"gpuMemory"`
}

type NodeType struct {
	ID          string          `json:"id"`
	Name        *string         `json:"name"`
	Price       *int            `json:"price"`
	Regions     []string        `json:"regions"`
	Description *string         `json:"description"`
	Provider    RuntimeProvider `json:"provider"`
	Specs       NodeSpecs       `json:"specs"`
}

type Node struct {
	ID       string          `json:"id"`
	TypeID   string          `json:"typeID"`
	Region   string          `json:"region"`
	KeyPair  SSHKey          `json:"sshKeyPair"`
	Status   SessionStatus   `json:"status"`
	Provider RuntimeProvider `json:"provider"`
}

type Project struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SSHKey struct {
	Name      string     `json:"name"`
	PublicKey *string    `json:"publicKey,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type Session struct {
	ID         uuid.UUID       `json:"id"`
	SSHKey     SSHKey          `json:"sshKey"`
	Status     SessionStatus   `json:"runtimeStatus"`
	NodeTypeID string          `json:"nodeTypeID"`
	Region     string          `json:"region"`
	Provider   RuntimeProvider `json:"provider"`
}

type ExecParams struct {
	Cmd   []string `json:"cmd"`
	Image []string `json:"containerImage"`
}
