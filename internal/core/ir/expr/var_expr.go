package expr

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

type Var struct {
	Rvalue
	Lvalue
	Name   string
	Type   lang.Type
	Global bool
}

func (v Var) String() string {
	return v.Name
}
func (v Var) GetType() lang.Type {
	return v.Type
}
