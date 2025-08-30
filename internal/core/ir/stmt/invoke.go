package stmt

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
)

// InvokeStmt invoke function or method and assign return value
type InvokeStmt struct {
	DefinitionStmt
	Lvalue     expr.Var
	Rvalue     expr.Invoke
	Container  *lang.PHPMethod
	LineNumber int
	Seq        int
}

func (i *InvokeStmt) GetLValue() expr.Lvalue {
	return &i.Lvalue
}

func (i *InvokeStmt) GetRValue() expr.Rvalue {
	return i.Rvalue
}

func (i *InvokeStmt) GetLineNumber() int {
	return i.LineNumber
}

func (i *InvokeStmt) String() string {
	if i.Lvalue.String() == "" {
		return i.Rvalue.String()
	}
	return i.Lvalue.String() + "=" + i.Rvalue.String()
}

func (i *InvokeStmt) GetInvokeExpr() expr.Invoke {
	return i.Rvalue.(expr.Invoke)
}

func (i *InvokeStmt) Accept(visitor StmtVisitor) {
	visitor.VisitInvokeStaticStmt(*i)
}

func (i *InvokeStmt) GetStmtSeq() int {
	return i.Seq
}

func (i *InvokeStmt) SetStmtSeq(i0 int) {
	i.Seq = i0
}
