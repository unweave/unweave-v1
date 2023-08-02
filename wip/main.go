package main

import (
	"github.com/unweave/unweave-v1/wip/conductor"
	"github.com/unweave/unweave-v1/wip/providers/local"
)

var (
	p1 = local.NewProvider("p1")
	p2 = local.NewProvider("p2")
)

func main() {
	// start the orchestrator
	// - one orchestrator for every provider

	conductor.RegisterProvider(p1)
	conductor.RegisterProvider(p2)

	// add api routes for exec creation and deletion
	// forward exec creation and deletion to exec service
	// exec service decides whether to use the Unweave Driver (kubernetes) or the Conductor (??)
	// The exec service should ask each type of driver to list availability
	// The driver that can serve the request the fastest should be used

}
