package mpscunboundedarrayqueue

import (
	"fmt"
	"testing"
	"time"
)

func TestData(t *testing.T) {
	mpsc := NewMpscUnboundedQueue[int](64)

	println(mpsc.producerMask)

	for i := 0; i < 3000; i++ {
		mpsc.Offer(i)
	}
	for i := 0; i < 3000; i++ {
		mpsc.RelaxedPoll()
	}

	println(mpsc.GetMask())

	time.Sleep(1 * time.Second)

	println(mpsc.producerIndex.Load())
	println(mpsc.consumerIndex.Load())

	//arr := mpsc.consumerBuffer
	//for i, p := range arr {
	//val := p
	//if val == nil {
	//fmt.Printf("index %d: nil\n", i)
	//} else {
	//fmt.Printf("index %d: %+v\n", i, *val.Load()) // dereference
	//}
	//}
	fmt.Println(mpsc.consumerBuffer)
	fmt.Println(mpsc.producerBuffer)
	//println(*mpsc.producerBuffer[0].Load())

	println(mpsc.producerLimit.Load())

	println(mpsc.producerMask)

}
