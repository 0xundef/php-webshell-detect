package lang

type MemberRef interface {
	String() string
	IsStatic() bool
}
