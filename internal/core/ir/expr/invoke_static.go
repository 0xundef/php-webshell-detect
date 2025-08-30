package expr

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

// InvokeStaticExp fun1()
type InvokeStaticExp struct {
	Invoke
	FunRef lang.MethodRef
	Args   []Var
	FunVar Var // function Name maybe stored in var
}

func (i *InvokeStaticExp) GetFunName() string {
	return i.FunRef.Name
}

func (i *InvokeStaticExp) String() string {
	var funName string
	if i.FunRef.Name != "" {
		funName = i.FunRef.String()
	}
	if i.FunVar.Name != "" {
		funName = i.FunVar.String()
	}
	return funName + "(" + i.GetArgList() + ")"
}
func (i *InvokeStaticExp) GetType() lang.Type {
	return i.FunRef.ReturnType
}
func (i *InvokeStaticExp) GetArgs() []Var {
	return i.Args
}

func (i *InvokeStaticExp) GetArgList() string {
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

func (i *InvokeStaticExp) GetContainer() *lang.PHPMethod {
	return nil
}
