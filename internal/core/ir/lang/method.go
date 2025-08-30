package lang

import (
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/ast"
	"strings"
)

type PHPMethod struct {
	ClassMember
	DeclaringClass *PHPClass
	Name           string

	ParamTypes []Type
	ReturnType Type
	ParamNames []string
	Source     ast.Vertex
}

func (method *PHPMethod) GetSignature() string {
	var className string
	if method.DeclaringClass == nil {
		className = "UN_KNOW_CLASS"
	} else {
		className = method.DeclaringClass.Name
	}
	args := strings.Join(method.ParamNames, ",")
	return className + "::" + method.Name + "(" + args + ")"
}

func CreatePHPMethod(declaringClass *PHPClass,
	name string,
	paramTypes []Type,
	returnType Type,
	paramNames []string) *PHPMethod {

	return &PHPMethod{
		DeclaringClass: declaringClass,
		Name:           name,
		ParamTypes:     paramTypes,
		ReturnType:     returnType,
		ParamNames:     paramNames,
	}
}
