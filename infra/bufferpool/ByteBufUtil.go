package bufferpool

import (
	"bytes"
	"sync"
)

var BufferPool sync.Pool

func init() {
	BufferPool = sync.Pool{
		New: func() interface{} {
			var buffer bytes.Buffer
			var bufferp *bytes.Buffer = &buffer
			return bufferp
		},
	}
}

func NewBufferWithLength(length int) *bytes.Buffer {
	buffer := BufferPool.Get().(*bytes.Buffer)

	buffer.Grow(length)

	return buffer
}
