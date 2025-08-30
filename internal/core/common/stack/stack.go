package stack

type Stack struct {
	buf  []interface{}
	size int
}

func NewStack(size int) IStackInterface {
	if size < 4 {
		size = 4
	}
	return &Stack{
		buf: make([]interface{}, 0, size),
	}
}

func (sk *Stack) Push(x interface{}) {
	if len(sk.buf) == sk.size {
		sk.buf = append(sk.buf, x)
	} else {
		sk.buf[sk.size] = x
	}
	sk.size++
}

func (sk *Stack) Pop() (interface{}, bool) {
	if sk.size == 0 {
		return nil, false
	}
	val := sk.buf[sk.size-1]
	sk.buf[sk.size-1] = nil //释放对变量的引用
	sk.size--
	return val, true
}

func (sk *Stack) Top() (interface{}, bool) {
	if sk.size == 0 {
		return nil, false
	}
	return sk.buf[sk.size-1], true
}

func (sk *Stack) IsEmpty() bool {
	return sk.size == 0
}

func (sk *Stack) Len() int {
	return sk.size
}

func (sk *Stack) Cap() int {
	return cap(sk.buf)
}
