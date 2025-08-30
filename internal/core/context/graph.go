package context

type Edge interface {
	GetSource() interface{}
	GetTarget() interface{}
}

type EdgeKind int

const (
	LOCAL EdgeKind = iota
	STORE
	LOAD
	PARAMETER_PASSING
	RETURN
)

func (kind EdgeKind) String() string {
	return [...]string{"local", "store", "Kind", "."}[kind]
}
