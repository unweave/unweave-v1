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
	setMinIfZero := func(val *int, min int) {
		if *val == 0 {
			*val = min
		}
	}

	setMaxIfZeroOrBelowMin := func(min, max *int) {
		if *max <= *min {
			*max = *min
		}
	}

	setDefaultGPU := func(val *string) {
		if *val == "" {
			*val = defaultGPURequest
		}
	}

	setDefaultGPU(&spec.GPU.Type)

	setMinIfZero(&spec.GPU.Count.Min, defaultMinGPUs)
	setMaxIfZeroOrBelowMin(&spec.GPU.Count.Min, &spec.GPU.Count.Max)

	setMinIfZero(&spec.CPU.Min, defaultMinCPU)
	setMaxIfZeroOrBelowMin(&spec.CPU.Min, &spec.CPU.Max)

	setMinIfZero(&spec.HDD.Min, defaultMinHDD)
	setMaxIfZeroOrBelowMin(&spec.HDD.Min, &spec.HDD.Max)
	return spec
}
