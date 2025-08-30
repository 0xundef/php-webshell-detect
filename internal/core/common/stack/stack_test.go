package stack

import (
	"fmt"
	"testing"
)

func TestNewStack(t *testing.T) {
	s := NewStack(10)
	size := 10
	for i := 0; i < size; i++ {
		s.Push(i)
	}
	if s.Len() != size {
		t.Fatal("len error")
	}
	t.Logf("cap=%v", s.Cap())
	for !s.IsEmpty() {
		if val, exist := s.Top(); exist {
			fmt.Println(val)
		}
		s.Pop()
	}
	if s.Len() != 0 {
		t.Fatal("empty error")
	}
}
