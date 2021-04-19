package go_ds

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	s := NewStack()
	for i := 0; i < 10; i++ {
		s.Push(i)
	}

	for j := 0; j < 11; j++ {
		fmt.Println(s.Pop())
	}
}
