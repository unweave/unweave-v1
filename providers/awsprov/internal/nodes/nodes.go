package nodes

import (
	"math"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/unweave/unweave/api/types"
)

type AwsNodeType struct {
	NodeID   ec2types.InstanceType
	NodeName string
	Cost     float64
	GPUMem   int
	GPUCount int
	CPUCount int
	CPUMem   float64
}

type nodeOptions []AwsNodeType

func (opts nodeOptions) atLeastGPUs(count int) nodeOptions {
	if len(opts) == 0 {
		return nil
	}

	var out nodeOptions

	for i := range opts {
		opt := opts[i]

		if opt.GPUCount >= count {
			out = append(out, opt)
		}
	}

	return out
}

func (opts nodeOptions) atLeastGPUMem(mem int) nodeOptions {
	if len(opts) == 0 {
		return nil
	}

	var out nodeOptions

	for i := range opts {
		opt := opts[i]

		if opt.GPUMem >= mem {
			out = append(out, opt)
		}
	}

	return out
}

func (opts nodeOptions) atLeastCPUCount(count int) nodeOptions {
	if len(opts) == 0 {
		return nil
	}

	var out nodeOptions

	for i := range opts {
		opt := opts[i]

		if opt.CPUCount >= count {
			out = append(out, opt)
		}
	}

	return out
}

func (opts nodeOptions) atLeastCPUMem(mem int) nodeOptions {
	if len(opts) == 0 {
		return nil
	}

	var out nodeOptions

	for i := range opts {
		opt := opts[i]

		if opt.CPUMem >= float64(mem) {
			out = append(out, opt)
		}
	}

	return out
}

func (opts nodeOptions) cheapest() *AwsNodeType {
	if len(opts) == 0 {
		return nil
	}

	var cheapest *AwsNodeType

	for idx := range opts {
		if cheapest == nil {
			cheapest = &opts[idx]
		}

		if opts[idx].Cost < cheapest.Cost {
			cheapest = &opts[idx]
		}
	}

	return cheapest
}

func BestFit(nodes []AwsNodeType, spec types.HardwareSpec) ec2types.InstanceType {
	opts := nodeOptions(nodes)

	node := opts.
		atLeastGPUs(spec.GPU.Count.Min).
		atLeastGPUMem(spec.GPU.RAM.Min).
		atLeastCPUCount(spec.CPU.Min).
		atLeastCPUMem(spec.RAM.Min).
		cheapest()

	if node == nil {
		// no nodes found
		return ""
	}

	return node.NodeID
}

func GeneralPurposeNodes() []AwsNodeType {
	return []AwsNodeType{
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeT3Nano, CPUCount: 2, CPUMem: 0.5},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeT3Micro, CPUCount: 2, CPUMem: 1},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeT3Small, CPUCount: 2, CPUMem: 2},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeT3Medium, CPUCount: 2, CPUMem: 4},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeM6iLarge, CPUCount: 2, CPUMem: 8},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeM6iXlarge, CPUCount: 4, CPUMem: 16},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeM6i2xlarge, CPUCount: 8, CPUMem: 32},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeM6i4xlarge, CPUCount: 16, CPUMem: 64},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeM6i8xlarge, CPUCount: 32, CPUMem: 128},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeM6i12xlarge, CPUCount: 48, CPUMem: 192},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeM6i16xlarge, CPUCount: 64, CPUMem: 256},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeM6i24xlarge, CPUCount: 96, CPUMem: 384},
		{NodeName: "General purpose", NodeID: ec2types.InstanceTypeM6i32xlarge, CPUCount: 128, CPUMem: 512},
	}
}

func GPUNodes() map[string][]AwsNodeType {
	return map[string][]AwsNodeType{
		"tesla_m60": {
			{NodeName: "Tesla M60", NodeID: ec2types.InstanceTypeG3sXlarge, Cost: 0, GPUCount: 1, GPUMem: 8, CPUCount: 4, CPUMem: 30.5},
			{NodeName: "Tesla M60", NodeID: ec2types.InstanceTypeG34xlarge, Cost: 0, GPUCount: 1, GPUMem: 8, CPUCount: 16, CPUMem: 122},
			{NodeName: "Tesla M60", NodeID: ec2types.InstanceTypeG38xlarge, Cost: 0, GPUCount: 1, GPUMem: 8, CPUCount: 32, CPUMem: 244},
			{NodeName: "Tesla M60", NodeID: ec2types.InstanceTypeG316xlarge, Cost: 0, GPUCount: 1, GPUMem: 8, CPUCount: 64, CPUMem: 488},
		},
		"t4_tensor": {
			{NodeName: "T4 Tensor", NodeID: ec2types.InstanceTypeG4dnXlarge, Cost: 0, GPUCount: 1, CPUCount: 4, CPUMem: 16, GPUMem: 16},
			{NodeName: "T4 Tensor", NodeID: ec2types.InstanceTypeG4dn2xlarge, Cost: 0, GPUCount: 1, CPUCount: 8, CPUMem: 32, GPUMem: 16},
			{NodeName: "T4 Tensor", NodeID: ec2types.InstanceTypeG4dn4xlarge, Cost: 0, GPUCount: 1, CPUCount: 16, CPUMem: 64, GPUMem: 16},
			{NodeName: "T4 Tensor", NodeID: ec2types.InstanceTypeG4dn8xlarge, Cost: 0, GPUCount: 1, CPUCount: 32, CPUMem: 128, GPUMem: 16},
			{NodeName: "T4 Tensor", NodeID: ec2types.InstanceTypeG4dn16xlarge, Cost: 0, GPUCount: 1, CPUCount: 64, CPUMem: 256, GPUMem: 16},
			{NodeName: "T4 Tensor", NodeID: ec2types.InstanceTypeG4dn12xlarge, Cost: 0, GPUCount: 4, CPUCount: 48, CPUMem: 192, GPUMem: 64},
		},
		"aws_inferentia": {
			{NodeName: "AWS Inferentia", NodeID: ec2types.InstanceTypeInf1Xlarge, Cost: 0, GPUCount: 1, CPUCount: 4, CPUMem: 8, GPUMem: 8},
			{NodeName: "AWS Inferentia", NodeID: ec2types.InstanceTypeInf12xlarge, Cost: 0, GPUCount: 1, CPUCount: 8, CPUMem: 16, GPUMem: 16},
			{NodeName: "AWS Inferentia", NodeID: ec2types.InstanceTypeInf16xlarge, Cost: 0, GPUCount: 4, CPUCount: 24, CPUMem: 48, GPUMem: 48},
			{NodeName: "AWS Inferentia", NodeID: ec2types.InstanceTypeInf124xlarge, Cost: 0, GPUCount: 16, CPUCount: 96, CPUMem: 192, GPUMem: 192},
		},
		"aws_inferentia2": {
			{NodeName: "AWS Inferentia2", NodeID: ec2types.InstanceTypeInf2Xlarge, Cost: 0, GPUCount: 1, GPUMem: 32, CPUCount: 4, CPUMem: 16},
			{NodeName: "AWS Inferentia2", NodeID: ec2types.InstanceTypeInf28xlarge, Cost: 0, GPUCount: 1, GPUMem: 32, CPUCount: 32, CPUMem: 128},
			{NodeName: "AWS Inferentia2", NodeID: ec2types.InstanceTypeInf224xlarge, Cost: 0, GPUCount: 6, GPUMem: 192, CPUCount: 96, CPUMem: 384},
			{NodeName: "AWS Inferentia2", NodeID: ec2types.InstanceTypeInf248xlarge, Cost: 0, GPUCount: 12, GPUMem: 384, CPUCount: 192, CPUMem: 768},
		},
		"aws_trainium": {
			{NodeName: "AWS Trainium", NodeID: ec2types.InstanceTypeTrn12xlarge, Cost: 0, GPUCount: 1, GPUMem: 32, CPUCount: 8, CPUMem: 32},
			{NodeName: "AWS Trainium", NodeID: ec2types.InstanceTypeTrn132xlarge, Cost: 0, GPUCount: 16, GPUMem: 512, CPUCount: 128, CPUMem: 512},
			{NodeName: "AWS Trainium", NodeID: ec2types.InstanceTypeTrn1n32xlarge, Cost: 0, GPUCount: 16, GPUMem: 512, CPUCount: 128, CPUMem: 512},
		},
		"gaudi": {
			{NodeName: "Gaudi", NodeID: ec2types.InstanceTypeDl124xlarge, Cost: 0, CPUCount: 96, GPUCount: 8, CPUMem: 768, GPUMem: 768},
		},
		"k80": {
			{NodeName: "K80", NodeID: ec2types.InstanceTypeP2Xlarge, Cost: 0, GPUCount: 1, CPUCount: 4, CPUMem: 61, GPUMem: 12},
			{NodeName: "K80", NodeID: ec2types.InstanceTypeP28xlarge, Cost: 0, GPUCount: 8, CPUCount: 32, CPUMem: 488, GPUMem: 96},
			{NodeName: "K80", NodeID: ec2types.InstanceTypeP216xlarge, Cost: 0, GPUCount: 16, CPUCount: 64, CPUMem: 732, GPUMem: 192},
		},
		"tesla_v100": {
			{NodeName: "Tesla V100", NodeID: ec2types.InstanceTypeP32xlarge, Cost: 0, GPUCount: 1, CPUCount: 8, CPUMem: 61, GPUMem: 16},
			{NodeName: "Tesla V100", NodeID: ec2types.InstanceTypeP38xlarge, Cost: 0, GPUCount: 4, CPUCount: 32, CPUMem: 244, GPUMem: 64},
			{NodeName: "Tesla V100", NodeID: ec2types.InstanceTypeP316xlarge, Cost: 0, GPUCount: 8, CPUCount: 64, CPUMem: 488, GPUMem: 128},
			{NodeName: "Tesla V100", NodeID: ec2types.InstanceTypeP3dn24xlarge, Cost: 0, GPUCount: 8, CPUCount: 96, CPUMem: 768, GPUMem: 256},
		},
		"a100": {
			{NodeName: "A100 Tensor", NodeID: ec2types.InstanceTypeP4d24xlarge, Cost: 0, GPUCount: 8, CPUCount: 96, CPUMem: 1152, GPUMem: 320},
		},
	}
}

func NodeType(spec types.HardwareSpec) ec2types.InstanceType {
	if spec.GPU.Type == "" {
		if spec.CPU.Type != "intel" {
			return ec2types.InstanceType(spec.CPU.Type)
		}

		return BestFit(GeneralPurposeNodes(), spec)
	}

	gpuNodeOptions := GPUNodes()[spec.GPU.Type]

	return BestFit(gpuNodeOptions, spec)
}

func CPUNodeTypes() types.NodeType {
	return toNodeType("CPU", "intel", "intel", "", GeneralPurposeNodes())
}

func toNodeType(nodeType string, nodeID string, cpuType string, gpuType string, cpuNodes []AwsNodeType) types.NodeType {
	smallest := cpuNodes[0]
	largest := cpuNodes[len(cpuNodes)-1]

	return types.NodeType{
		Type:     nodeType,
		ID:       nodeID,
		Name:     &smallest.NodeName,
		Price:    nil,
		Regions:  []string{},
		Provider: types.AWSProvider,
		Specs: types.HardwareSpec{
			GPU: types.GPU{
				Type: gpuType,
				Count: types.HardwareRequestRange{
					Min: smallest.GPUCount,
					Max: largest.GPUCount,
				},
				RAM: types.HardwareRequestRange{
					Min: smallest.GPUMem,
					Max: largest.GPUMem,
				},
			},
			CPU: types.CPU{
				Type: cpuType,
				HardwareRequestRange: types.HardwareRequestRange{
					Min: smallest.CPUCount,
					Max: largest.CPUCount,
				},
			},
			RAM: types.HardwareRequestRange{
				Min: int(math.Ceil(smallest.CPUMem)),
				Max: int(math.Ceil(largest.CPUMem)),
			},
			HDD: types.HardwareRequestRange{
				Min: 0,
				Max: 2000,
			},
		},
	}
}

func ToNodeTypesGPU(gpuNodes map[string][]AwsNodeType) []types.NodeType {
	out := make([]types.NodeType, 0, len(gpuNodes))

	for k, v := range gpuNodes {
		out = append(out, toNodeType("GPU", k, "", k, v))
	}

	return out
}
