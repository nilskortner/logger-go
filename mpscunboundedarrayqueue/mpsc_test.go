package mpscunboundedarrayqueue

import (
	"sync"
	"testing"
)

func TestData(t *testing.T) {
	mpsc := NewMpscUnboundedQueue[int](1024)

	//println(mpsc.producerMask)

	counter := 100_000

	var wait sync.WaitGroup
	wait.Add(counter)

	for i := 0; i < counter; i++ {
		go func() {
			mpsc.Offer(i)
			wait.Done()
		}()
	}

	wait.Wait()

	//println(mpsc.GetMask())

	//time.Sleep(1 * time.Second)

	//println(mpsc.producerIndex.Load())
	//println(mpsc.consumerIndex.Load())

	//arr := mpsc.consumerBuffer
	//for i, p := range arr {
	//val := p
	//if val == nil {
	//fmt.Printf("index %d: nil\n", i)
	//} else {
	//fmt.Printf("index %d: %+v\n", i, *val.Load()) // dereference
	//}
	//}
	//fmt.Println(mpsc.consumerBuffer)
	//fmt.Println(mpsc.producerBuffer)
	//println(*mpsc.producerBuffer[0].Load())

	//println(mpsc.producerLimit.Load())

	//println(mpsc.producerMask)

}
