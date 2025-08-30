package ir

type KIND int

const (
	UNDEF KIND = iota
	CONST_INT
	CONST_STR
	NAC
)

func (k KIND) String() string {
	return [...]string{"UNDEF", "CONST_INT", "CONST_STR", "NAC"}[k]
}

type Value struct {
	Kind  KIND
	Value string
}

func (v Value) IsConst() bool {
	return v.Kind == CONST_INT || v.Kind == CONST_STR
}
