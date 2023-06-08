package main

import (
	"fmt"
)

func NewRingBuffer(inCh, outCh chan int) *ringBuffer {
	return &ringBuffer{
		inCh:  inCh,
		outCh: outCh,
	}
}

type ringBuffer struct {
	inCh  chan int
	outCh chan int
}

func (r *ringBuffer) Run() {
	defer close(r.outCh)
	for inChanVal := range r.inCh {
		select {
		case r.outCh <- inChanVal: // если в буффере канала есть место -> записываем сразу туда
		default: // если канал заполнен -> выччитываем старую запись и записываем новую
			oldVal := <-r.outCh
			fmt.Printf("Old value %d\n", oldVal)
			r.outCh <- inChanVal
		}
	}
}

func main() {
	max := 100
	inCh := make(chan int, max)
	outCh := make(chan int, 10)

	for i := 0; i < max; i++ {
		inCh <- i
	}

	rb := NewRingBuffer(inCh, outCh)
	close(inCh)
	rb.Run()

	resSlice := make([]int, 0)
	for res := range outCh {
		resSlice = append(resSlice, res)
	}
	fmt.Println(resSlice)
}
