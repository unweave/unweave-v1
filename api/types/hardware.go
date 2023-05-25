package types

type HardwareRequestRange struct {
	Min int `json:"min,omitempty"`
	Max int `json:"max,omitempty"`
}

type GPU struct {
	Count HardwareRequestRange `json:"count"`
	Type  string               `json:"type,omitempty"`
	RAM   HardwareRequestRange `json:"ram,omitempty"`
}

type HardwareSpec struct {
	GPU     GPU                  `json:"gpu"`
	CPU     HardwareRequestRange `json:"cpu"`
	RAM     HardwareRequestRange `json:"ram"`
	Storage HardwareRequestRange `json:"storage"`
}
