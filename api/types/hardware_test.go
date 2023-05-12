package types

import (
	"fmt"
	"testing"
)

func TestParseRange(t *testing.T) {
	tests := []struct {
		input    string
		expected HardwareRequestRange
		err      error
	}{
		{"", HardwareRequestRange{}, nil},
		{"1", HardwareRequestRange{1, 1}, nil},
		{"1-3", HardwareRequestRange{1, 3}, nil},
		{
			"1-",
			HardwareRequestRange{},
			fmt.Errorf("invalid range - must be of the form <number> or <number>-<number>"),
		},
		{
			"-3",
			HardwareRequestRange{},
			fmt.Errorf("invalid range - must be of the form <number> or <number>-<number>"),
		},
		{
			"3-1",
			HardwareRequestRange{},
			fmt.Errorf("invalid range - max range should be greater than or equal to min range"),
		},
	}

	for _, test := range tests {
		result, err := parseRange(test.input)
		if result != test.expected {
			t.Errorf("Expected %v, got %v for input %q", test.expected, result, test.input)
		}
		if err != nil && test.err != nil && err.Error() != test.err.Error() {
			t.Errorf("Expected error %v, got %v for input %q", test.err, err, test.input)
		}
	}
}

func TestParseGPU(t *testing.T) {
	tests := []struct {
		input    string
		expected GPU
		err      error
	}{
		{"", GPU{}, nil},
		{"1", GPU{HardwareRequestRange{1, 1}, ""}, nil},
		{"1-2", GPU{HardwareRequestRange{1, 2}, ""}, nil},
		{"1-2_nvidia", GPU{HardwareRequestRange{1, 2}, "nvidia"}, nil},
		{"1_nvidia", GPU{HardwareRequestRange{1, 1}, "nvidia"}, nil},
		{"_nvidia", GPU{HardwareRequestRange{}, "nvidia"}, nil},
		{
			"1-2-",
			GPU{},
			fmt.Errorf("invalid GPU spec - must be of the form G<number>_<type>"),
		},
		{
			"1-2-nvidia",
			GPU{},
			fmt.Errorf("invalid GPU spec - must be of the form G<number>_<type>"),
		},
	}

	for _, test := range tests {
		result, err := parseGPU(test.input)
		if result != test.expected {
			t.Errorf("Expected %v, got %v for input %q", test.expected, result, test.input)
		}
		if err != nil && test.err != nil && err.Error() != test.err.Error() {
			t.Errorf("Expected error %v, got %v for input %q", test.err, err, test.input)
		}
	}
}

func TestHardwareSpecParse(t *testing.T) {
	tests := []struct {
		input    string
		expected *HardwareSpec
		err      bool
	}{
		{"", &HardwareSpec{}, false},
		{
			"G1,C1,R1,S1",
			&HardwareSpec{
				GPU{HardwareRequestRange{1, 1}, ""},
				HardwareRequestRange{1, 1},
				HardwareRequestRange{1, 1},
				HardwareRequestRange{1, 1},
			},
			false,
		},
		{
			"G1_NVIDIA,C1-2,R1-2,S1-2",
			&HardwareSpec{
				GPU{HardwareRequestRange{1, 1}, "nvidia"},
				HardwareRequestRange{1, 2},
				HardwareRequestRange{1, 2},
				HardwareRequestRange{1, 2},
			},
			false,
		},
		{
			"G,C,R,S",
			&HardwareSpec{},
			true,
		},
		// Test case for positional syntax.
		{
			"1,1,1,1",
			&HardwareSpec{
				GPU{HardwareRequestRange{1, 1}, ""},
				HardwareRequestRange{1, 1},
				HardwareRequestRange{1, 1},
				HardwareRequestRange{1, 1},
			},
			false,
		},
		{
			"1-3,2,4-9",
			&HardwareSpec{
				GPU{HardwareRequestRange{1, 3}, ""},
				HardwareRequestRange{2, 2},
				HardwareRequestRange{4, 9},
				HardwareRequestRange{0, 0},
			},
			false,
		},
		{
			"a100",
			&HardwareSpec{
				GPU{HardwareRequestRange{1, 1}, "a100"},
				HardwareRequestRange{0, 0},
				HardwareRequestRange{0, 0},
				HardwareRequestRange{0, 0},
			},
			false,
		},
		{
			"g2_a40,23,34",
			&HardwareSpec{},
			true,
		},
		{
			"0,10",
			&HardwareSpec{
				GPU: GPU{
					Count: HardwareRequestRange{
						Min: 0,
						Max: 0,
					},
					Type: "",
				},
				CPU: HardwareRequestRange{
					Min: 10,
					Max: 10,
				},
			},
			false,
		},
	}

	for _, test := range tests {
		hwSpec := &HardwareSpec{}
		err := hwSpec.Parse(test.input)
		if *hwSpec != *test.expected {
			t.Errorf("Expected %v, got %v for input %q", test.expected, hwSpec, test.input)
		}
		if err != nil && test.err && err == nil {
			t.Errorf("Expected error for input %q", test.input)
		}
	}
}
