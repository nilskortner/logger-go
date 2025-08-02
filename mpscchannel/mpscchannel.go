package mpscchannel

import (
	"runtime"
	"sync/atomic"
)

type MpscUnboundedChannel[T comparable] struct {
	head    *ChannelView[T]
	tail    *ChannelView[T]
	data    atomic.Pointer[chan T]
	size    int
	counter atomic.Int64
}

func NewMpscChannel[T comparable](size int) *MpscUnboundedChannel[T] {
	if size < 2 {
		return nil
	}
	data := make(chan T, size)
	view := &ChannelView[T]{
		data: data,
	}
	if size%2 != 0 {
		size += 1
	}
	mpsc := &MpscUnboundedChannel[T]{
		head: view,
		tail: view,
		size: size,
	}
	mpsc.data.Store(&data)
	return mpsc
}

func (m *MpscUnboundedChannel[T]) Offer(value T) bool {
start:
	if m.counter.Load()%1 == 1 {
		runtime.Gosched()
		goto start
	} else if m.counter.Load() < int64(m.size*2) {
		select {
		case *m.data.Load() <- value:
			m.counter.Add(2)
			return true
		default:
			//sync.Cond wait
			goto start
		}
	} else if m.counter.Load() == int64(m.size*2) {
		m.resize(int64(m.size * 2))
		goto start
	}
	goto start
}

func (m *MpscUnboundedChannel[T]) resize(pSize int64) {
	if m.counter.CompareAndSwap(pSize, pSize+1) {
		defer m.counter.Store(0)
		data := make(chan T, m.size)
		m.head.next = &ChannelView[T]{
			data: data,
		}
		m.data.Store(&data)
		m.head = m.head.next
		println("resize: ")
	}
}

func (m *MpscUnboundedChannel[T]) RelaxedPoll() (T, bool) {
start:
	select {
	case value := <-m.tail.data:
		m.counter.Add(-2)
		return value, true
	default:
		if m.tail.next == nil {
			var zeroValue T
			return zeroValue, false
		} else {
			m.tail = m.tail.next
			goto start
		}
	}
}

func (m *MpscUnboundedChannel[T]) TestGetSize() int {
	return m.size
}

func (m *MpscUnboundedChannel[T]) TestGetChannel() int {
	return len(*m.data.Load())
}

func (m *MpscUnboundedChannel[T]) TestGetCounter() int {
	return int(m.counter.Load())
}
