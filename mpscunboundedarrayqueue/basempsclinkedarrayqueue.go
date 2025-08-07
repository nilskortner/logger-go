package mpscunboundedarrayqueue

import (
	"fmt"
	"loggergo/mathsupport"
	"loggergo/mpscunboundedarrayqueue/util"
	"sync/atomic"
	"time"
)

var CONTINUE_TO_P_INDEX_CAS int = 0
var RETRY int = 1
var QUEUE_FULL int = 2
var QUEUE_RESIZE int = 3

type Buffer[T any] struct {
	data []*atomic.Pointer[T]
	next *Buffer[T]
}

type BaseMpscLinkedArrayQueueProducerFields struct {
	_             [64]byte
	producerIndex atomic.Int64
	_             [64]byte
}

func (pf *BaseMpscLinkedArrayQueueProducerFields) lvProducerIndex() int64 {
	return (&pf.producerIndex).Load()
}

func (pf *BaseMpscLinkedArrayQueueProducerFields) TestingLvProducerIndex() int64 {
	return (&pf.producerIndex).Load()
}

func (pf *BaseMpscLinkedArrayQueueProducerFields) soProducerIndex(newValue int64) {
	(&pf.producerIndex).Store(newValue)
}

func (pf *BaseMpscLinkedArrayQueueProducerFields) casProducerIndex(expect, newValue int64) bool {
	return (&pf.producerIndex).CompareAndSwap(expect, newValue)
}

type BaseMpscLinkedArrayQueueConsumerFields[T any] struct {
	_              [64]byte
	consumerIndex  atomic.Int64
	consumerMask   int64
	consumerBuffer []*atomic.Pointer[T]
	_              [64]byte
}

func (cf *BaseMpscLinkedArrayQueueConsumerFields[T]) GetCBuffer() []*atomic.Pointer[T] {
	return cf.consumerBuffer
}

func (cf *BaseMpscLinkedArrayQueueConsumerFields[T]) lvConsumerIndex() int64 {
	return (&cf.consumerIndex).Load()
}

func (cf *BaseMpscLinkedArrayQueueConsumerFields[T]) TestingLvConsumerIndex() int64 {
	return (&cf.consumerIndex).Load()
}

func (cf *BaseMpscLinkedArrayQueueConsumerFields[T]) soConsumerIndex(newValue int64) {
	(&cf.consumerIndex).Store(newValue)
}

type BaseMpscLinkedArrayQueueColdProducerFields[T any] struct {
	_              [64]byte
	producerLimit  atomic.Int64
	producerMask   int64
	producerBuffer []*atomic.Pointer[T]
	_              [64]byte
}

func (cpf *BaseMpscLinkedArrayQueueColdProducerFields[T]) GetMask() int64 {
	return cpf.producerMask
}

func (cpf *BaseMpscLinkedArrayQueueColdProducerFields[T]) GetBuffer() []*atomic.Pointer[T] {
	return cpf.producerBuffer
}

func (cpf *BaseMpscLinkedArrayQueueColdProducerFields[T]) TestingLvProducerLimit() int64 {
	return (&cpf.producerLimit).Load()
}

func (cpf *BaseMpscLinkedArrayQueueColdProducerFields[T]) lvProducerLimit() int64 {
	return (&cpf.producerLimit).Load()
}

func (cpf *BaseMpscLinkedArrayQueueColdProducerFields[T]) casProducerLimit(expect, newValue int64) bool {
	return (&cpf.producerLimit).CompareAndSwap(expect, newValue)
}

func (cpf *BaseMpscLinkedArrayQueueColdProducerFields[T]) soProducerLimit(newValue int64) {
	(&cpf.producerLimit).Store(newValue)
}

type BaseMpscLinkedArrayQueue[T any] struct {
	*BaseMpscLinkedArrayQueueProducerFields
	*BaseMpscLinkedArrayQueueConsumerFields[T]
	*BaseMpscLinkedArrayQueueColdProducerFields[T]
	Head     *Buffer[T]
	Tail     *Buffer[T]
	Capacity int64
}

func NewBaseMpscLinkedArrayQueue[T any](initialCapacity int) *BaseMpscLinkedArrayQueue[T] {
	_, err := util.CheckGreaterThanOrEqual(initialCapacity, 2, "initialCapacity")
	if err != nil {
		fmt.Println(err)
	}

	p2capacity := mathsupport.RoundToPowerOfTwo(initialCapacity)

	mask := int64(p2capacity-1) << 1

	capacity := int64(p2capacity + 1)

	buffer := make([]*atomic.Pointer[T], capacity)

	firstBuffer := &Buffer[T]{data: buffer}

	bmlaq := &BaseMpscLinkedArrayQueue[T]{
		Capacity:                               capacity,
		Head:                                   firstBuffer,
		Tail:                                   firstBuffer,
		BaseMpscLinkedArrayQueueProducerFields: &BaseMpscLinkedArrayQueueProducerFields{},
		BaseMpscLinkedArrayQueueConsumerFields: &BaseMpscLinkedArrayQueueConsumerFields[T]{
			consumerMask:   mask,
			consumerBuffer: buffer,
		},
		BaseMpscLinkedArrayQueueColdProducerFields: &BaseMpscLinkedArrayQueueColdProducerFields[T]{
			producerMask:   mask,
			producerBuffer: buffer,
		},
	}

	bmlaq.BaseMpscLinkedArrayQueueColdProducerFields.soProducerLimit(mask)

	return bmlaq
}

func (b *BaseMpscLinkedArrayQueue[T]) TestingGetMovingBuffer() *Buffer[T] {
	return b.Head
}

func (b *BaseMpscLinkedArrayQueue[T]) TestingGetMovingBufferData(buf *Buffer[T]) []*atomic.Pointer[T] {
	return buf.data
}

func (b *BaseMpscLinkedArrayQueue[T]) TestingGetMultiMovingBufferData() ([]*atomic.Pointer[T], []*atomic.Pointer[T], []*atomic.Pointer[T], []*atomic.Pointer[T], []*atomic.Pointer[T]) {
	return b.Head.data, b.Head.next.data, b.Head.next.next.data, b.Head.next.next.next.data, b.Head.next.next.next.next.data
}

func (b *BaseMpscLinkedArrayQueue[T]) TestingGetConsumerMask() int64 {
	return b.consumerMask
}

func (b *BaseMpscLinkedArrayQueue[T]) Offer(e T) bool {
	p := &e

	var mask int64
	var buffer []*atomic.Pointer[T]
	var pIndex int64

	/// Fix For Cas Contention. Not in Tests
	//attempts := 0

	for {
		producerLimit := b.lvProducerLimit()
		pIndex = b.lvProducerIndex()

		if (pIndex & 1) == 1 {
			/// Fix For Cas Contention
			//attempts++
			//backoff := time.Duration(rand.Intn(1<<uint(min(attempts, 10)))) * time.Microsecond
			//time.Sleep(backoff)
			///
			continue
		}

		mask = b.producerMask
		buffer = b.producerBuffer

		if producerLimit < pIndex {
			result := b.offerSlowPath(mask, pIndex, producerLimit)
			switch result {
			case CONTINUE_TO_P_INDEX_CAS:
				break
			case RETRY:
				continue
			case QUEUE_FULL:
				return false
			case QUEUE_RESIZE:
				b.resize(buffer, pIndex, p)
				return true
			}
		}

		if b.casProducerIndex(pIndex, pIndex+2) {
			break
		}
	}
	//INDEX visible before ELEMENT
	offset := pIndex & mask
	//println(pIndex & mask)
	soRefElement(buffer, offset, p)
	return true
}

func (b *BaseMpscLinkedArrayQueue[T]) RelaxedPoll() (T, bool) {
	var zeroValue T

	buffer := b.consumerBuffer
	cIndex := b.lvConsumerIndex()
	mask := b.consumerMask

	offset := cIndex & mask
	e := lvRefElement[T](buffer, offset)
	if e == nil {
		if buffer[b.Capacity-1] != nil {
			nextBuffer := b.nextBuffer()
			valuePointer := b.newBufferPoll(nextBuffer, cIndex)
			if valuePointer != nil {
				return *valuePointer, true
			} else {
				return zeroValue, false
			}
		}
		return zeroValue, false
	}
	soRefElement(buffer, offset, nil)
	b.soConsumerIndex(cIndex + 2)
	return *e, true
}

func (b *BaseMpscLinkedArrayQueue[T]) offerSlowPath(mask, pIndex, producerLimit int64) int {
	cIndex := b.lvConsumerIndex()
	bufferCapacity := getCurrentBufferCapacity(mask)

	if cIndex+bufferCapacity > pIndex {
		if !b.casProducerLimit(producerLimit, cIndex+bufferCapacity) {
			// 1 = retry from top
			return RETRY
		} else {
			// 0 = continue to pIndex CAS
			return CONTINUE_TO_P_INDEX_CAS
		}
		// full and cannot grow
	} else if availableInQueue(pIndex, cIndex) <= 0 {
		// 2 = Queue full. offer should return false
		return QUEUE_FULL
		// grab index for resize -> set lower bit
	} else if b.casProducerIndex(pIndex, pIndex+1) {
		// 3 = trigger a resize
		return QUEUE_RESIZE
	} else {
		// failed resize attempt, retry from top
		return RETRY
	}
}

func (b *BaseMpscLinkedArrayQueue[T]) resize(oldBuffer []*atomic.Pointer[T], pIndex int64, p *T) {
	if p == nil {
		panic("no clear value defined in func resize()")
	}

	//TIMER
	timer := time.Now()
	//

	// make new JUMP Value Pointer
	var jumpVal T
	jump := &jumpVal

	newBufferLength := b.Capacity

	//
	// Risk of Running out of Memory
	//
	newBuffer := make([]*atomic.Pointer[T], newBufferLength)

	b.producerBuffer = newBuffer
	newMask := (newBufferLength - 2) << 1
	b.producerMask = newMask

	var offsetInOld int64 = b.Capacity - 1
	offsetInNew := pIndex & newMask

	soRefElement(newBuffer, offsetInNew, p)
	b.appendNext(newBuffer)

	// ASSERT code
	cIndex := b.lvConsumerIndex()
	availableInQueue := availableInQueue(pIndex, cIndex)
	util.CheckPositive(availableInQueue, "availableInQueue")

	// Invalidate racing CASs
	// We mever set the limit beyond the bounds of a buffer
	b.soProducerLimit(pIndex + mathsupport.MinInt64(newMask, availableInQueue))

	// INDEX visible before ELEMENT, consistent with consumer expectation

	// make resize visible to consumer
	soRefElement(oldBuffer, offsetInOld<<1, jump)

	// TIMER
	timing := time.Now().UnixMilli() - timer.UnixMilli()
	//if timing > 1 {
	println(timing)
	//}
	///

	// make resize visible to the other producers
	b.soProducerIndex(pIndex + 2)

}

func (b *BaseMpscLinkedArrayQueue[T]) nextBuffer() []*atomic.Pointer[T] {
	b.Head = b.Head.next
	var nextBuffer []*atomic.Pointer[T] = b.Head.data

	b.consumerBuffer = nextBuffer
	b.consumerMask = int64(len(nextBuffer)-2) << 1
	return nextBuffer
}

func (b *BaseMpscLinkedArrayQueue[T]) newBufferPoll(nextBuffer []*atomic.Pointer[T], cIndex int64) *T {
	offset := cIndex & b.consumerMask
	n := lvRefElement[T](nextBuffer, offset)
	if n == nil {
		return n
	}
	soRefElement(nextBuffer, offset, nil)
	b.soConsumerIndex(cIndex + 2)
	return n
}

func (b *BaseMpscLinkedArrayQueue[T]) appendNext(nextBuffer []*atomic.Pointer[T]) {
	b.Tail.next = (&Buffer[T]{data: nextBuffer})
	b.Tail = b.Tail.next
}

func lvRefElement[T any](buffer []*atomic.Pointer[T], index int64) *T {
	index = index >> 1
	if buffer[index] == nil {
		return nil
	}
	return buffer[index].Load()
}

func soRefElement[T any](buffer []*atomic.Pointer[T], index int64, value *T) {
	index = index >> 1
	if value == nil {
		buffer[index] = nil
		return
	}
	if buffer[index] == nil {
		buffer[index] = &atomic.Pointer[T]{}
	}
	buffer[index].Store(value)
}
