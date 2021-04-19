package go_ds

import (
	"fmt"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	q := NewQueue()
	go func(q *Queue) {
		for {
			fmt.Println(111111)
			fmt.Println(q.Pop())
		}
	}(q)

	for i := 0; i < 10; i++ {
		q.Push(1)
		time.Sleep(time.Second * 2)
	}
}
