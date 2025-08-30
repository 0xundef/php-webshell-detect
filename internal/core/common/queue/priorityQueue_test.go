package queue

import (
	"testing"
)

type integer int

func (i integer) GetPriority() int {
	return int(i)
}

func (i integer) GetHashCode() int {
	return int(i)
}

func TestNewPriorityQueue(t *testing.T) {
	count := 100
	heap := NewPriorityQueue(count)
	for i := 0; i < count; i++ {
		heap.Push(integer(i))
	}
	if !heap.Exist(integer(20)) {
		t.Fatal("Contain error")
	}
	if !heap.Remove(integer(20)) {
		t.Fatal("remove error")
	}
	heap.Push(integer(20))
	for heap.Len() > 0 {
		l := heap.Len()
		val := heap.Pop()
		if val.GetHashCode() != count-l {
			t.Fatal("pop error", val.GetHashCode(), count-l)
		}
	}
	for i := 0; i < count; i++ {
		heap.Push(integer(i))
	}
	list := heap.GetArray()
	t.Log(list)
	heap.Pop()
	list = heap.GetArray()
	t.Log(list)
	//heap.Push(integer(20))
	//heap.Push(integer(20))
	list = heap.GetArray()
	t.Log(list)
	heap.Pop()
	for heap.Len() > 0 {
		t.Log(heap.Pop())
	}
	if !heap.Empty() {
		t.Fatal("clear error")
	}
	heap.Clear()
	if !heap.Empty() {
		t.Fatal("clear error")
	}

}

//
//func Test(t *testing.T)  {
//	a := New(5)
//	a.Put("1")
//	a.Put("1")
//	a.Put("1")
//	a.Put("1")
//	fmt.Println(a.Len())
//	a.Peek()
//	fmt.Println(a.Len())
//	a.Get(1)
//	fmt.Println(a.Len())
//	a.Dispose()
//}
