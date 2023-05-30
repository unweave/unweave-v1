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

func SetSpecDefaultValues(spec HardwareSpec) HardwareSpec {
	const defaultMinValueRequest = 1
	const defaultGPURequest = "rtx_4000"

	setMinIfZero := func(val *int) {
		if *val == 0 {
			*val = defaultMinValueRequest
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

	setMinIfZero(&spec.GPU.Count.Min)
	setMaxIfZeroOrBelowMin(&spec.GPU.Count.Min, &spec.GPU.Count.Max)

	setMinIfZero(&spec.GPU.RAM.Min)
	setMaxIfZeroOrBelowMin(&spec.GPU.RAM.Min, &spec.GPU.RAM.Max)

	setMinIfZero(&spec.CPU.Min)
	setMaxIfZeroOrBelowMin(&spec.CPU.Min, &spec.CPU.Max)

	setMinIfZero(&spec.RAM.Min)
	setMaxIfZeroOrBelowMin(&spec.RAM.Min, &spec.RAM.Max)

	setMinIfZero(&spec.HDD.Min)
	setMaxIfZeroOrBelowMin(&spec.HDD.Min, &spec.HDD.Max)

	return spec
}
