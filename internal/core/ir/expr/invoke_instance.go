package expr

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

// InvokeInstanceExp $obj->method()
type InvokeInstanceExp struct {
	Invoke
	MethodRef lang.MethodRef
	Args      []Var
	Base      Var
}

func (i *InvokeInstanceExp) GetFunName() string {
	return i.MethodRef.Name
}

func (i *InvokeInstanceExp) String() string {
	return i.Base.String() + "->" + i.MethodRef.Name + "(" + i.GetArgList() + ")"
}
func (i *InvokeInstanceExp) GetType() lang.Type {
	return i.MethodRef.ReturnType
}
func (i *InvokeInstanceExp) GetContainer() *lang.PHPMethod {
	return nil
}
func (i *InvokeInstanceExp) GetArgs() []Var {
	return i.Args
}

func (i *InvokeInstanceExp) GetArgList() string {
	var ret string
	for i0 := 0; i0 < len(i.Args); i0++ {
		if i0 != 0 {
			ret = ret + "," + i.Args[i0].Name
		} else {
			ret = i.Args[i0].Name
		}
	}
	return ret
}
