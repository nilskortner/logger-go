package mathsupport

import "math/bits"

var maxpow int = 1 << 30

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
	if x > maxpow {
		return -1
	}
	if x < 0 {
		return -1
	}

	return 1 << (32 - bits.LeadingZeros32(uint32(x-1)))
}
