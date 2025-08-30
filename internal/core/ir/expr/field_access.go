package expr

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

// FieldAccess $obj->name
type FieldAccess struct {
	Lvalue
	Rvalue
	Base     Var
	FieldRef lang.FieldRef
}

func (f *FieldAccess) String() string {
	return f.Base.String() + "->" + f.FieldRef.Name
}

func (f *FieldAccess) GetType() lang.Type {
	return nil
}
