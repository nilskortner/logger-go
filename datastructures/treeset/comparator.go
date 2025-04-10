package treeset

import "loggergo/logging/core/appender/file/logfile"

type Comparator[T comparable] func(x, y T) int

func LogComparator(x, y logfile.LogFile) int {
	switch {
	case x.GetIndex() > y.GetIndex():
		return 1
	case x.GetIndex() < y.GetIndex():
		return -1
	default:
		return 0
	}
}

func StringComparator(x, y string) int {
	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}
