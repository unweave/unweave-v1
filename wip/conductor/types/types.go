package types

type Network struct {
	Ports map[int]int
}

type Spec struct {
	CPU  int
	RAM  int
	GPUs []string
	Network
}

type NodeState string

const (
	NodeInitializing NodeState = "initializing"
	NodeRunning      NodeState = "running"
	NodeIdle         NodeState = "idle"
	NodeStopped      NodeState = "stopped"
	NodeError        NodeState = "error"
)

type ContainerState string

const (
	ContainerInitializing ContainerState = "initializing"
	ContainerRunning      ContainerState = "running"
	ContainerStopped      ContainerState = "stopped"
	ContainerError        ContainerState = "error"
	ContainerExited       ContainerState = "exited"
)
