package lang

func FromInt8(b1, b2, b3, b4 int8) int {
	return int(uint(b1)<<24 | uint(b2)<<16 | uint(b3)<<8 | uint(b4))
}
