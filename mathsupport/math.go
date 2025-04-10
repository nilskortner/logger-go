package mathsupport

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// MinInt64 returns the smaller of two int64 values.
// It compares a and b, and returns the smaller one.
func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func RoundToPowerOfTwo(x int) int {
	if x <= 0 {
		return -1
	}
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	if ^uint(0)>>32 != 0 {
		x |= x >> 32
	}
	return x + 1
}
