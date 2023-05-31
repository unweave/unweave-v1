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
	GPU GPU                  `json:"gpu"`
	CPU HardwareRequestRange `json:"cpu"`
	RAM HardwareRequestRange `json:"ram"`
	HDD HardwareRequestRange `json:"hdd"`
}

const (
	defaultMinCPU     = 1
	defaultMinHDD     = 4
	defaultMinGPUs    = 1
	defaultGPURequest = "rtx_4000"
)

func SetSpecDefaultValues(spec HardwareSpec) HardwareSpec {
	spec.GPU.Type = setDefaultGPU(spec.GPU.Type)

	spec.GPU.Count.Min = setMinIfZero(spec.GPU.Count.Min, defaultMinGPUs)
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

func setDefaultGPU(val string) string {
	if val == "" {
		val = defaultGPURequest
	}

	return val
}
