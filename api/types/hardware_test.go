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
		{"1_3", HardwareRequestRange{1, 3}, nil},
		{
			"1_",
			HardwareRequestRange{},
			fmt.Errorf("invalid range - must be of the form <number> or <number>_<number>"),
		},
		{
			"_3",
			HardwareRequestRange{},
			fmt.Errorf("invalid range - must be of the form <number> or <number>_<number>"),
		},
		{
			"3_1",
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
		{"1_2", GPU{HardwareRequestRange{1, 2}, ""}, nil},
		{"1_2(NVIDIA)", GPU{HardwareRequestRange{1, 2}, "NVIDIA"}, nil},
		{"1(NVIDIA)", GPU{HardwareRequestRange{1, 1}, "NVIDIA"}, nil},
		{"(NVIDIA)", GPU{HardwareRequestRange{}, "NVIDIA"}, nil},
		{
			"1_2_",
			GPU{},
			fmt.Errorf("invalid GPU spec - must be of the form G<number>(<type>)"),
		},
		{
			"1_2_(NVIDIA)",
			GPU{},
			fmt.Errorf("invalid GPU spec - must be of the form G<number>(<type>)"),
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
		err      error
	}{
		{"", &HardwareSpec{}, nil},
		{
			"G1,C1,R1,S1",
			&HardwareSpec{
				GPU{HardwareRequestRange{1, 1}, ""},
				HardwareRequestRange{1, 1},
				HardwareRequestRange{1, 1},
				HardwareRequestRange{1, 1},
			},
			nil,
		},
		{
			"G1(NVIDIA),C1_2,R1_2,S1_2",
			&HardwareSpec{
				GPU{HardwareRequestRange{1, 1}, "NVIDIA"},
				HardwareRequestRange{1, 2},
				HardwareRequestRange{1, 2},
				HardwareRequestRange{1, 2},
			},
			nil,
		},
		{
			"G,C,R,S",
			&HardwareSpec{},
			fmt.Errorf("invalid range - must be of the form <number> or <number>_<number>"),
		},
	}

	for _, test := range tests {
		hwSpec := &HardwareSpec{}
		err := hwSpec.Parse(test.input)
		if *hwSpec != *test.expected {
			t.Errorf("Expected %v, got %v for input %q", test.expected, hwSpec, test.input)
		}
		if err != nil && test.err != nil && err.Error() != test.err.Error() {
			t.Errorf("Expected error %v, got %v for input %q", test.err, err, test.input)
		}
	}
}
