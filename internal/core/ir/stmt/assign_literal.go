package stmt

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
)

// AssignLiteralStmt $a="a" $b=5
type AssignLiteralStmt struct {
	DefinitionStmt
	Lvalue     expr.Lvalue
	Rvalue     expr.Literal
	LineNumber int
	Seq        int
}

func (a *AssignLiteralStmt) GetLValue() expr.Lvalue {
	return a.Lvalue
}

func (a *AssignLiteralStmt) GetRValue() expr.Rvalue {
	return a.Rvalue
}

func (a *AssignLiteralStmt) GetLineNumber() int {
	return a.LineNumber
}

func (a *AssignLiteralStmt) String() string {
	return a.Lvalue.String() + "=" + a.Rvalue.String()
}

func (a *AssignLiteralStmt) Accept(visitor StmtVisitor) {
	visitor.VisitAssignLiteralStmt(*a)
}

func (a *AssignLiteralStmt) GetStmtSeq() int {
	return a.Seq
}
func (a *AssignLiteralStmt) SetStmtSeq(i int) {
	a.Seq = i
}
