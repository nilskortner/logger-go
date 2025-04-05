package util

import "fmt"

func CheckGreaterThanOrEqual(n, expected int, name string) (int, error) {
	if n < expected {
		return 0, fmt.Errorf("%s: %d (expected: >= %d)", name, n, expected)
	}

	return n, nil
}

func CheckPositive(n int64, name string) (int64, error) {
	if n <= 0 {
		return 0, fmt.Errorf(name+": %d (expected > 0)", n)
	}

	return n, nil
}
