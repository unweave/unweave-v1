package conductor

import (
	"errors"
	"fmt"
	"time"

	"github.com/unweave/unweave-v1/wip/conductor/node"
	"github.com/unweave/unweave-v1/wip/conductor/types"
	"github.com/unweave/unweave-v1/wip/conductor/volume"
)

// Each user has their own orchestrator
// Orc maintains a pool of nodes to serve the user request
// It uses the providers registered when creating the orchestrator

type Store struct {
	Volume volume.Store
}

type orchestrator struct {
	provider     Provider
	nodes        map[string]node.Node
	containers   map[string]container
	volumes      map[string]volume.Volume
	nodePoolSpec types.Spec
	poolSize     int
}

func (o *orchestrator) newNode(spec types.Spec) {

}

func (o *orchestrator) terminateNode(nodeID string) {

}

func (o *orchestrator) start() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			if len(o.nodes) < o.poolSize {
				o.newNode(o.nodePoolSpec)
			}

			for _, _ = range o.nodes {

			}
		}
	}
}

var (
	ErrorNoNode        error = errors.New("no node available")
	ErrorOutOfCapacity error = errors.New("out of capacity")
)

func (o *orchestrator) match(spec types.Spec) (nodeID string, err error) {
	for _, n := range o.nodes {
		if n.State == types.NodeIdle {
			if n.Spec.CPU >= spec.CPU &&
				n.Spec.RAM >= spec.RAM &&
				len(n.Spec.GPUs) >= len(spec.GPUs) {
				return n.ID, nil
			}
		}
	}
	return "", ErrorNoNode
}

func (o *orchestrator) handle() {

}

func (o *orchestrator) assign(spec types.Spec) (containerID, nodeID string, err error) {
	// find a node that matches the spec
	// if no node matches, create a new node
	// if no node can be created, return error

	nodeID, err = o.match(spec)
	if err == nil {

		//cid := o.provider.CreateContainer(spec)
		cid := "123"
		o.containers[cid] = container{
			ID:     cid,
			NodeID: nodeID,
			Spec:   spec,
		}

		return cid, nodeID, nil
	}

	if err == ErrorNoNode {
		o.newNode(spec)
		return o.assign(spec)
	}

	return "", "", fmt.Errorf("failed to assign spec: %w", err)

}

func new(provider Provider) *orchestrator {
	return &orchestrator{
		provider: provider,
	}
}

type ContainerCreateConfig struct {
	WorkingDir string
	Cmd        []string
	SSHKeys    []string
	Volumes    map[string]string
}

var orcMap map[string]*orchestrator
var numWorkers int

func Init(providers map[string]Provider) {

	for p := range providers {
		// start a goroutine for each provider

		for i := 0; i < numWorkers; i++ {
			o := new(p)
			go o.start()
		}
	}
}

func RegisterProvider(provider Provider) {
	o := new(provider)

	for i := 0; i < numWorkers; i++ {
		go o.handle()
	}
	orcMap[provider.ID()] = o
}
