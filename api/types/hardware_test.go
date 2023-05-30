package types

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestSetDefaultValues(t *testing.T) {
	g := Goblin(t)

	const defaultMinValue = 1
	const defaultGPU = "rtx_4000"

	g.Describe("SetDefaultValues", func() {
		g.It("should set default values for unset fields", func() {
			g.Describe("when all fields are unset", func() {
				existingSpec := HardwareSpec{}
				expectedSpec := HardwareSpec{
					GPU: GPU{
						Count: HardwareRequestRange{
							Min: defaultMinValue,
							Max: defaultMinValue,
						},
						RAM: HardwareRequestRange{
							Min: defaultMinValue,
							Max: defaultMinValue,
						},
						Type: defaultGPU,
					},
					CPU: HardwareRequestRange{
						Min: defaultMinValue,
						Max: defaultMinValue,
					},
					RAM: HardwareRequestRange{
						Min: defaultMinValue,
						Max: defaultMinValue,
					},
					HDD: HardwareRequestRange{
						Min: defaultMinValue,
						Max: defaultMinValue,
					},
				}

				actualSpec := SetSpecDefaultValues(existingSpec)
				g.Assert(actualSpec).Eql(expectedSpec)
			})

			g.Describe("when some fields have specific values set", func() {
				existingSpec := HardwareSpec{
					GPU: GPU{
						Count: HardwareRequestRange{
							Min: 2,
						},
						RAM: HardwareRequestRange{
							Max: 4,
						},
						Type: "test",
					},
					CPU: HardwareRequestRange{
						Max: 8,
					},
					RAM: HardwareRequestRange{
						Min: 16,
					},
					HDD: HardwareRequestRange{
						Min: 32,
						Max: 16,
					},
				}
				expectedSpec := HardwareSpec{
					GPU: GPU{
						Count: HardwareRequestRange{
							Min: 2,
							Max: 2,
						},
						RAM: HardwareRequestRange{
							Min: defaultMinValue,
							Max: 4,
						},
						Type: "test",
					},
					CPU: HardwareRequestRange{
						Min: defaultMinValue,
						Max: 8,
					},
					RAM: HardwareRequestRange{
						Min: 16,
						Max: 16,
					},
					HDD: HardwareRequestRange{
						Min: 32,
						Max: 32,
					},
				}

				actualSpec := SetSpecDefaultValues(existingSpec)
				g.Assert(actualSpec).Eql(expectedSpec)
			})
		})

		g.It("should maintain existing values for already set fields", func() {
			g.Describe("when all fields are already set", func() {
				existingSpec := HardwareSpec{
					GPU: GPU{
						Count: HardwareRequestRange{
							Min: 2,
							Max: 4,
						},
						RAM: HardwareRequestRange{
							Min: 8,
							Max: 16,
						},
					},
					CPU: HardwareRequestRange{
						Min: 4,
						Max: 8,
					},
					RAM: HardwareRequestRange{
						Min: 16,
						Max: 32,
					},
					HDD: HardwareRequestRange{
						Min: 32,
						Max: 64,
					},
				}
				expectedSpec := HardwareSpec{
					GPU: GPU{
						Count: HardwareRequestRange{
							Min: 2,
							Max: 4,
						},
						RAM: HardwareRequestRange{
							Min: 8,
							Max: 16,
						},
						Type: defaultGPU,
					},
					CPU: HardwareRequestRange{
						Min: 4,
						Max: 8,
					},
					RAM: HardwareRequestRange{
						Min: 16,
						Max: 32,
					},
					HDD: HardwareRequestRange{
						Min: 32,
						Max: 64,
					},
				}

				actualSpec := SetSpecDefaultValues(existingSpec)
				g.Assert(actualSpec).Eql(expectedSpec)
			})
		})
	})
}
