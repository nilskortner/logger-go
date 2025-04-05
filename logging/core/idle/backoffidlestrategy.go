package idle

import (
	"gurms/internal/supportpkgs/mathsupport"
	"runtime"
	"strconv"
	"time"
)

const WORKING = 0
const SPINNING = 1
const YIELDING = 2
const PARKING = 3

type BackoffIdleStrategy struct {
	maxSpins        int64
	maxYields       int64
	minParkPeriodNs int64
	maxParkPeriodNs int64
	state           int
	value           int64
}

func NewBackoffIdleStrategy(maxSpins, maxYields, minParkPeriodNs, maxParkPeriodNs int64) *BackoffIdleStrategy {
	if minParkPeriodNs < 1 || maxParkPeriodNs < minParkPeriodNs {
		panic(
			"The minimum park period (" +
				strconv.FormatInt(minParkPeriodNs, 10) +
				") is less than 1, " +
				"and the maximum park (" +
				strconv.FormatInt(maxParkPeriodNs, 10) +
				") period is less than the minimum park period")
	}
	return &BackoffIdleStrategy{
		maxSpins:        maxSpins,
		maxYields:       maxYields,
		minParkPeriodNs: minParkPeriodNs,
		maxParkPeriodNs: maxParkPeriodNs,
		state:           WORKING,
	}
}

func (b *BackoffIdleStrategy) Idle() {
	switch b.state {
	case WORKING:
		b.value = 0
		b.state = SPINNING
		// fallthrough

	case SPINNING:
		if b.value+1 <= b.maxSpins {
			break
		}
		b.value = 0
		b.state = YIELDING

	case YIELDING:
		if b.value+1 <= b.maxYields {
			runtime.Gosched()
			break
		}
		b.value = b.minParkPeriodNs
		b.state = PARKING
		// fallthrough

	case PARKING:
		time.Sleep(time.Duration(b.value) * time.Nanosecond)
		b.value = mathsupport.MinInt64(b.value<<1, b.maxParkPeriodNs)
	}
}

func (b *BackoffIdleStrategy) Reset() {
	b.state = WORKING
}
