package types

func Inty(i int) *int             { return &i }
func Booly(b bool) *bool          { return &b }
func Stringy(s string) *string    { return &s }
func Float64y(f float64) *float64 { return &f }

func IntInv(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func BoolInv(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func StringInv(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Float64Inv(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
