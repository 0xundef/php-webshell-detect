package expr

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

// ArrayAccess $_GET["a"]
type ArrayAccess struct {
	Lvalue
	Rvalue
	Base  Var
	Index int
}

func (f *ArrayAccess) String() string {
	return f.Base.String() + "[*]"
}

func (f *ArrayAccess) GetType() lang.Type {
	return nil
}
