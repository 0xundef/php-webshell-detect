package ir

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/ast"
	"strconv"
)

func MakeTmpVar(n ast.Vertex, t lang.Type) expr.Var {
	return expr.Var{
		Name: "%temp" + strconv.Itoa(n.GetPosition().StartLine) + "_" + strconv.Itoa(n.GetPosition().StartPos),
		Type: t,
	}
}

func MakeConstVar(n ast.Vertex, t lang.Type) expr.Var {
	return expr.Var{
		Name: "%const" + t.String() + strconv.Itoa(n.GetPosition().StartLine) + "_" + strconv.Itoa(n.GetPosition().StartPos),
		Type: t,
	}
}

func MakeVar(name string, t lang.Type, global bool) expr.Var {
	return expr.Var{
		Name:   name,
		Type:   t,
		Global: global,
	}
}

func ConvertString(node ast.Vertex) string {
	var ret string
	if id, ok := node.(*ast.Identifier); ok {
		ret = string(id.Value)
	}
	return ret
}
