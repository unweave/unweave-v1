package types

type HardwareRequestRange struct {
	Min int `json:"min,omitempty"`
	Max int `json:"max,omitempty"`
}

type GPU struct {
	Type  string               `json:"type,omitempty"`
	Count HardwareRequestRange `json:"count"`
	RAM   HardwareRequestRange `json:"ram,omitempty"`
}

type HardwareSpec struct {
	GPU GPU                  `json:"gpu"`
	CPU HardwareRequestRange `json:"cpu"`
	RAM HardwareRequestRange `json:"ram"`
	HDD HardwareRequestRange `json:"hdd"`
}

const (
	defaultMinCPU = 4
	defaultMinHDD = 50
)

func SetSpecDefaultValues(spec HardwareSpec) HardwareSpec {
	spec.GPU.Count.Min = setDefaultMinGPUCount(spec.GPU)
	spec.GPU.Count.Max = setMaxIfZeroOrBelowMin(spec.GPU.Count.Min, spec.GPU.Count.Max)

	spec.CPU.Min = setMinIfZero(spec.CPU.Min, defaultMinCPU)
	spec.CPU.Max = setMaxIfZeroOrBelowMin(spec.CPU.Min, spec.CPU.Max)

	spec.HDD.Min = setMinIfZero(spec.HDD.Min, defaultMinHDD)
	spec.HDD.Max = setMaxIfZeroOrBelowMin(spec.HDD.Min, spec.HDD.Max)

	return spec
}

func setMinIfZero(val int, min int) int {
	if val == 0 {
		return min
	}
	return val
}

func setMaxIfZeroOrBelowMin(min, max int) int {
	if max <= min {
		return min
	}

	return max
}

func setDefaultMinGPUCount(gpu GPU) int {
	if gpu.Count.Min != 0 {
		return gpu.Count.Min
	}

	// User specified a GPU and didn't specify a min count. Use 1 as the default
	if gpu.Type != "" {
		return 1
	}

	// If no GPU type and no min count assume the user does not want a GPU
	return 0
}

type NodeMetadataV1 struct {
	ID             string           `json:"id"`
	TypeID         string           `json:"typeID"`
	Price          int              `json:"price"`
	VCPUs          int              `json:"vcpus"`
	Memory         int              `json:"memory"`
	HDD            int              `json:"hdd"`
	GpuType        string           `json:"gpuType"`
	GPUCount       int              `json:"gpuCount"`
	GPUMemory      int              `json:"gpuMemory"`
	ConnectionInfo ConnectionInfoV1 `json:"connection_info"`
}

func (m NodeMetadataV1) GetHardwareSpec() HardwareSpec {
	return HardwareSpec{
		GPU: GPU{
			Count: HardwareRequestRange{
				Min: m.GPUCount,
				Max: m.GPUCount,
			},
			Type: m.GpuType,
			RAM: HardwareRequestRange{
				Min: m.GPUMemory,
				Max: m.GPUMemory,
			},
		},
		CPU: HardwareRequestRange{
			Min: m.VCPUs,
			Max: m.VCPUs,
		},
		RAM: HardwareRequestRange{
			Min: m.Memory,
			Max: m.Memory,
		},
		HDD: HardwareRequestRange{
			Min: m.HDD,
			Max: m.HDD,
		},
	}
}

func DBNodeMetadataFromNode(node Node) NodeMetadataV1 {
	n := NodeMetadataV1{
		ID:        node.ID,
		TypeID:    node.TypeID,
		Price:     node.Price,
		VCPUs:     node.Specs.CPU.Min,
		Memory:    node.Specs.RAM.Min,
		HDD:       node.Specs.HDD.Min,
		GpuType:   node.Specs.GPU.Type,
		GPUCount:  node.Specs.GPU.Count.Min,
		GPUMemory: node.Specs.GPU.RAM.Min,

		ConnectionInfo: ConnectionInfoV1{
			Version: 1,
			Host:    node.Host,
			Port:    node.Port,
			User:    node.User,
		},
	}
	return n
}

type ConnectionInfoV1 struct {
	Version int    `json:"version"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	User    string `json:"user"`
}

func (c ConnectionInfoV1) GetConnectionInfo() *ConnectionInfo {
	return &ConnectionInfo{
		Host: c.Host,
		Port: c.Port,
		User: c.User,
	}
}
