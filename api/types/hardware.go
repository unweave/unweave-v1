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

	re := regexp.MustCompile(`^(\d+)(?:-(\d+))?$`)
	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return HardwareRequestRange{}, fmt.Errorf("invalid range - must be of the form <number> or <number>-<number>")
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

	// Match string that specific min gpus, max gpus, type or any combination of the three
	// Examples:
	// 1-2_nvidia, 1_nvidia, 1-2, 1, _nvidia
	re := regexp.MustCompile(`^(\d+)?(?:-(\d+))?(?:_(.+))?$`)
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

	return GPU{}, fmt.Errorf("invalid GPU spec - must be of the form G<number>_<type>")
}

func (h *HardwareSpec) String() string {
	var parts []string

	gpu := ""
	if h.GPU.Count.Min != 0 {
		gpu = fmt.Sprintf("%d", h.GPU.Count.Min)
	}
	if h.GPU.Count.Max != 0 {
		gpu = fmt.Sprintf("%s-%d", gpu, h.GPU.Count.Max)
	}
	if h.GPU.Type != "" {
		gpu = fmt.Sprintf("%s_%s", gpu, h.GPU.Type)
	}
	if gpu != "" {
		gpu = fmt.Sprintf("G%s", gpu)
	}

	cpu := ""
	if h.CPU.Min != 0 {
		cpu = fmt.Sprintf("%d", h.CPU.Min)
	}
	if h.CPU.Max != 0 {
		cpu = fmt.Sprintf("%s-%d", cpu, h.CPU.Max)
	}
	if cpu != "" {
		cpu = fmt.Sprintf("C%s", cpu)
	}

	ram := ""
	if h.RAM.Min != 0 {
		ram = fmt.Sprintf("%d", h.RAM.Min)
	}
	if h.RAM.Max != 0 {
		ram = fmt.Sprintf("%s-%d", ram, h.RAM.Max)
	}
	if ram != "" {
		ram = fmt.Sprintf("R%s", ram)
	}

	storage := ""
	if h.Storage.Min != 0 {
		storage = fmt.Sprintf("%d", h.Storage.Min)
	}
	if h.Storage.Max != 0 {
		storage = fmt.Sprintf("%s-%d", storage, h.Storage.Max)
	}
	if storage != "" {
		storage = fmt.Sprintf("S%s", storage)
	}

	parts = append(parts, gpu, cpu, ram, storage)
	return strings.Join(parts, ",")
}

func (h *HardwareSpec) Parse(spec string) error {
	spec = strings.ToLower(spec)
	parts := strings.Split(spec, ",")

	errmsg := fmt.Errorf("invalid hardware spec - must be of the form G<num_gpus>_<gpu_type>,C<num_cpus>,R<ram_gb>,S<storage_gb>")
	if len(parts) > 4 {
		return errmsg
	}

	validPositionalArgs := true

	// Match a string that starts with C, R, S, c, r, s and may or may not have a range
	// Examples:
	// C1, C1-2, c1, c1-2
	re := regexp.MustCompile(`^[crsCRS](\d*)(?:-(\d+))?$`)
	// Match a string that starts with G, g and may or may not have a range and a gpu type
	// Examples:
	// G1_nvidia, G1-2_nvidia, g1_nvidia, g1-2_nvidia
	gpuRe := regexp.MustCompile(`^[gG](\d*)?(?:-(\d+))?(?:_([a-zA-Z]+))?$`)

	for _, p := range parts {
		if len(p) == 0 {
			return errmsg
		}

		if match := re.FindStringSubmatch(p); match != nil {
			validPositionalArgs = false
			break
		}
		if match := gpuRe.FindStringSubmatch(p); match != nil {
			validPositionalArgs = false
			break
		}
	}

	var err error

	if validPositionalArgs {
		if len(parts) > 0 {
			// If the range is not present, append `1_` to the string to make it compatible with the gpu parser
			p1 := regexp.MustCompile(`^\d+-\d+_`)
			p2 := regexp.MustCompile(`^\d+_`)
			p3 := regexp.MustCompile(`^\d+`)
			p4 := regexp.MustCompile(`^_`)

			p := parts[0]
			if !p1.MatchString(p) && !p2.MatchString(p) && !p3.MatchString(p) && !p4.MatchString(p) {
				p = fmt.Sprintf("1_%s", p)
			}
			if h.GPU, err = parseGPU(p); err != nil {
				return err
			}
		}
		if len(parts) > 1 {
			if h.CPU, err = parseRange(parts[1]); err != nil {
				return err
			}
		}
		if len(parts) > 2 {
			if h.RAM, err = parseRange(parts[2]); err != nil {
				return err
			}
		}
		if len(parts) > 3 {
			if h.Storage, err = parseRange(parts[3]); err != nil {
				return err
			}
		}
		return nil
	}

	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		switch p[0] {
		case 'g':
			if h.GPU, err = parseGPU(strings.TrimLeft(p, "g")); err != nil {
				return err
			}
		case 'c':
			if h.CPU, err = parseRange(strings.TrimLeft(p, "c")); err != nil {
				return err
			}
		case 'r':
			if h.RAM, err = parseRange(strings.TrimLeft(p, "r")); err != nil {
				return err
			}
		case 's':
			if h.Storage, err = parseRange(strings.TrimLeft(p, "s")); err != nil {
				return err
			}
		}
	}
	return nil
}
