package ir

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
	"strconv"
	"strings"
)

// IR the basic statement container, it could be static function and class defined method
type IR interface {
	GetStmts() []stmt.Stmt
	GetStmt(i int) stmt.Stmt
	AddStmt(stmt stmt.Stmt)
	GetParamAt(i int) expr.Var
	AddParam(v expr.Var)
	GetParams() []expr.Var
	GetRetVars() []expr.Var

	AddReturn(v expr.Var)
	GetReturnAt(i int) expr.Var
	GetReturns() []expr.Var
	String() string
	GetThis() expr.Var
	GetContainerMethod() *lang.PHPMethod
	GetContainerClass() *lang.PHPClass
}

type DefaultIR struct {
	IR
	Method  *lang.PHPMethod
	Stmts   []stmt.Stmt
	Params  []expr.Var
	RetVars []expr.Var
	This    expr.Var
}

func (i *DefaultIR) GetParams() []expr.Var {
	return i.Params
}

func (i *DefaultIR) GetRetVars() []expr.Var {
	return i.RetVars
}

func (i *DefaultIR) GetContainerMethod() *lang.PHPMethod {
	return i.Method
}
func (i *DefaultIR) GetContainerClass() *lang.PHPClass {
	return i.Method.DeclaringClass
}
func (i *DefaultIR) GetStmts() []stmt.Stmt {
	return i.Stmts
}

func (i *DefaultIR) GetStmt(i0 int) stmt.Stmt {
	return i.Stmts[i0]
}

func (i *DefaultIR) AddStmt(stmt stmt.Stmt) {
	i.Stmts = append(i.Stmts, stmt)
}

func (i *DefaultIR) GetParamAt(i0 int) expr.Var {
	return i.Params[i0]
}

func (i *DefaultIR) AddReturn(v expr.Var) {
	i.RetVars = append(i.RetVars, v)
}

func (i *DefaultIR) AddParam(v expr.Var) {
	i.Params = append(i.Params, v)
}

func (i *DefaultIR) GetReturnAt(i0 int) expr.Var {
	return i.RetVars[i0]
}
func (i *DefaultIR) GetReturns() []expr.Var {
	return i.RetVars
}
func (i *DefaultIR) GetThis() expr.Var {
	return i.This
}
func (i *DefaultIR) String() string {
	var stmts []string
	stmts = append(stmts, i.Method.GetSignature())
	for _, s := range i.Stmts {
		stmts = append(stmts, strconv.Itoa(s.GetLineNumber())+" "+strconv.Itoa(s.GetStmtSeq())+"	"+s.String())
	}
	return strings.Join(stmts, "\n")
}
