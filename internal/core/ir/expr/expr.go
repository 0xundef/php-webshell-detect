package expr

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

type Expr interface {
	String() string
	GetType() lang.Type
}

type Rvalue interface {
	Expr
}

type Lvalue interface {
	Expr
}

type PrimitiveType int

const (
	INT PrimitiveType = iota
	CHAR
	BOOLEAN
	BYTE
	LONG
	FLOAT
	DOUBLE
	SHORT
)

func (t PrimitiveType) String() string {
	return [...]string{"int", "char", "boolean", "byte", "long", "float", "double", "short"}[t]
}
