package types

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestSetDefaultValues(t *testing.T) {
	g := Goblin(t)

	const defaultGPU = "rtx_4000"

	g.Describe("SetDefaultValues", func() {
		g.It("should set default values for unset fields", func() {
			g.Describe("when all fields are unset", func() {
				existingSpec := HardwareSpec{}
				expectedSpec := HardwareSpec{
					GPU: GPU{
						Count: HardwareRequestRange{
							Min: defaultMinGPUs,
							Max: defaultMinGPUs,
						},

						Type: defaultGPU,
					},
					CPU: HardwareRequestRange{
						Min: defaultMinCPU,
						Max: defaultMinCPU,
					},

					HDD: HardwareRequestRange{
						Min: defaultMinHDD,
						Max: defaultMinHDD,
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
						Type: "test",
					},
					CPU: HardwareRequestRange{
						Max: 8,
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
						Type: "test",
					},
					CPU: HardwareRequestRange{
						Min: defaultMinCPU,
						Max: 8,
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
					},
					CPU: HardwareRequestRange{
						Min: 4,
						Max: 8,
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
						Type: defaultGPU,
					},
					CPU: HardwareRequestRange{
						Min: 4,
						Max: 8,
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
