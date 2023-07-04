//nolint:paralleltest
package nodes_test

import (
	"testing"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/providers/awsprov/internal/nodes"
)

func TestNodesBestFit(t *testing.T) {
	t.Parallel()

	t.Run("CPU: should match 1 vCPU, 1GB", bestFitCPU(1, 1, ec2types.InstanceTypeT3Micro))
	t.Run("CPU: should match 10 vCPU, 1GB", bestFitCPU(10, 1, ec2types.InstanceTypeM6i4xlarge))
	t.Run("CPU: should match 1 vCPU, 100GB", bestFitCPU(1, 300, ec2types.InstanceTypeM6i24xlarge))

	t.Run(
		"GPU: should match 1x tesla_m60, 50x CPUs",
		bestFitGPU("tesla_m60", 1, 4, 50, 100, ec2types.InstanceTypeG316xlarge),
	)
	t.Run(
		"GPU: should match 5x tesla_v100, 10xCPUs, 500 GB CPU Ram",
		bestFitGPU("tesla_v100", 5, 4, 50, 500, ec2types.InstanceTypeP3dn24xlarge),
	)
}

func bestFitCPU(cpuCount, ram int, want ec2types.InstanceType) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		instance := nodes.BestFit(nodes.GeneralPurposeNodes(), cpuSpec(cpuCount, ram))
		assert.Equal(t, want, instance)
	}
}

func cpuSpec(cpuCount, ram int) types.HardwareSpec {
	return types.HardwareSpec{
		GPU: types.GPU{},
		CPU: types.CPU{
			Type: "x86_64",
			HardwareRequestRange: types.HardwareRequestRange{
				Min: cpuCount,
				Max: cpuCount,
			},
		},
		RAM: types.HardwareRequestRange{
			Min: ram,
			Max: ram,
		},
	}
}

func bestFitGPU(gpuType string, gpuCount, gpuMem, cpuCount, ram int, want ec2types.InstanceType) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		instance := nodes.BestFit(nodes.GPUNodes()[gpuType], gpuSpec(gpuType, gpuCount, gpuMem, cpuCount, ram))
		assert.Equal(t, want, instance)
	}
}

func gpuSpec(gpuType string, gpuCount, gpuRAM, cpuCount, ram int) types.HardwareSpec {
	return types.HardwareSpec{
		GPU: types.GPU{
			Type: gpuType,
			Count: types.HardwareRequestRange{
				Min: gpuCount,
				Max: gpuCount,
			},
			RAM: types.HardwareRequestRange{
				Min: gpuRAM,
				Max: gpuRAM,
			},
		},
		CPU: types.CPU{
			Type: "x86_64",
			HardwareRequestRange: types.HardwareRequestRange{
				Min: cpuCount,
				Max: cpuCount,
			},
		},
		RAM: types.HardwareRequestRange{
			Min: ram,
			Max: ram,
		},
	}
}
