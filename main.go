package main

import (
	"loggergo/mpscunboundedarrayqueue"
	"sync"
	"time"
)

func main() {

	t := time.Now().UnixMilli()

	counter := 1000_000

	queue := mpscunboundedarrayqueue.NewMpscUnboundedQueue[string](1024)

	var wait sync.WaitGroup
	wait.Add(counter)

	//factory.Loggerfactory(false, "test", nodetype.SERVICE, logging.NewLoggingProperties())
	//testlogger := factory.GetLogger("testlogger")

	for i := 0; i < counter; i++ {
		go func() {
			//testlogger.Fatal("fatal error")
			//time.Sleep(500 * time.Millisecond)
			queue.Offer("fatal error")
			wait.Done()
		}()
	}

	wait.Wait()

	println(time.Now().UnixMilli() - t)
}
