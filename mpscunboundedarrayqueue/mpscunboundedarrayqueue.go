package mpscunboundedarrayqueue

import (
	"math"
)

type MpscUnboundedArrayQueue[T any] struct {
	*BaseMpscLinkedArrayQueue[T]
}

func NewMpscUnboundedQueue[T any](chunkSize int) *MpscUnboundedArrayQueue[T] {
	return &MpscUnboundedArrayQueue[T]{
		BaseMpscLinkedArrayQueue: NewBaseMpscLinkedArrayQueue[T](chunkSize),
	}
}

func availableInQueue(pIndex, cIndex int64) int64 {
	_, _ = pIndex, cIndex // unused for this implementation
	return math.MaxInt64
}

func getCurrentBufferCapacity(mask int64) int64 {
	return mask
}
