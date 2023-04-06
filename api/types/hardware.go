package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type HardwareRequestRange struct {
	Min int
	Max int
}

type GPU struct {
	Count HardwareRequestRange
	Type  string
}

type HardwareSpec struct {
	GPU     GPU
	CPU     HardwareRequestRange
	RAM     HardwareRequestRange
	Storage HardwareRequestRange
}

func parseRange(s string) (HardwareRequestRange, error) {
	if s == "" {
		return HardwareRequestRange{}, nil
	}

	re := regexp.MustCompile(`^(\d+)(?:_(\d+))?$`)
	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return HardwareRequestRange{}, fmt.Errorf("invalid range - must be of the form <number> or <number>_<number>")
	}

	min, _ := strconv.Atoi(match[1])
	if match[2] == "" {
		return HardwareRequestRange{Min: min, Max: min}, nil
	}

	max, _ := strconv.Atoi(match[2])
	if max < min {
		return HardwareRequestRange{}, fmt.Errorf("invalid range - max range should be greater than or equal to min range")
	}
	return HardwareRequestRange{Min: min, Max: max}, nil
}

func parseGPU(s string) (GPU, error) {
	if s == "" {
		return GPU{}, nil
	}

	re := regexp.MustCompile(`^(?:(\d+)(?:_(\d+))?)?(?:\((\w+)\))?$`)
	match := re.FindStringSubmatch(s)

	if len(match) > 0 && (match[1] != "" || match[3] != "") {
		gpu := GPU{}

		if match[1] != "" {
			count, _ := strconv.Atoi(match[1])
			gpu.Count = HardwareRequestRange{Min: count, Max: count}
		}
		if match[2] != "" {
			maxCount, _ := strconv.Atoi(match[2])
			gpu.Count.Max = maxCount
		}
		if match[3] != "" {
			gpuType := match[3]
			gpu.Type = gpuType
		}
		return gpu, nil
	}

	return GPU{}, fmt.Errorf("invalid GPU spec - must be of the form G<number>(<type>)")
}

func (h *HardwareSpec) String() string {
	var parts []string

	gpu := ""
	if h.GPU.Count.Min != 0 {
		gpu = fmt.Sprintf("%d", h.GPU.Count.Min)
	}
	if h.GPU.Count.Max != 0 {
		gpu = fmt.Sprintf("%s_%d", gpu, h.GPU.Count.Max)
	}
	if h.GPU.Type != "" {
		gpu = fmt.Sprintf("%s(%s)", gpu, h.GPU.Type)
	}
	if gpu != "" {
		gpu = fmt.Sprintf("G%s", gpu)
	}

	cpu := ""
	if h.CPU.Min != 0 {
		cpu = fmt.Sprintf("%d", h.CPU.Min)
	}
	if h.CPU.Max != 0 {
		cpu = fmt.Sprintf("%s_%d", cpu, h.CPU.Max)
	}
	if cpu != "" {
		cpu = fmt.Sprintf("C%s", cpu)
	}

	ram := ""
	if h.RAM.Min != 0 {
		ram = fmt.Sprintf("%d", h.RAM.Min)
	}
	if h.RAM.Max != 0 {
		ram = fmt.Sprintf("%s_%d", ram, h.RAM.Max)
	}
	if ram != "" {
		ram = fmt.Sprintf("R%s", ram)
	}

	storage := ""
	if h.Storage.Min != 0 {
		storage = fmt.Sprintf("%d", h.Storage.Min)
	}
	if h.Storage.Max != 0 {
		storage = fmt.Sprintf("%s_%d", storage, h.Storage.Max)
	}
	if storage != "" {
		storage = fmt.Sprintf("S%s", storage)
	}

	return strings.Join(parts, ",")
}

func (h *HardwareSpec) Parse(spec string) error {
	spec = strings.ToUpper(spec)
	parts := strings.Split(spec, ",")

	var err error
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		switch p[0] {
		case 'G':
			if h.GPU, err = parseGPU(strings.TrimLeft(p, "G")); err != nil {
				return err
			}
		case 'C':
			if h.CPU, err = parseRange(strings.TrimLeft(p, "C")); err != nil {
				return err
			}
		case 'R':
			if h.RAM, err = parseRange(strings.TrimLeft(p, "R")); err != nil {
				return err
			}
		case 'S':
			if h.Storage, err = parseRange(strings.TrimLeft(p, "S")); err != nil {
				return err
			}
		}
	}
	return nil
}
