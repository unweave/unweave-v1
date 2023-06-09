package conductor

type container struct {
	ID     string
	NodeID string
	Status ContainerStatus
	Spec   Spec
}

func (c container) Logs() {

}

// watch container lifecycle
//
