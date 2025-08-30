package stmt

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
)

// AssignVarStmt var copy: $a = $b
type AssignVarStmt struct {
	DefinitionStmt
	Lvalue     expr.Var
	Rvalue     expr.Var
	LineNumber int
	Seq        int
}

func (a *AssignVarStmt) GetLValue() expr.Lvalue {
	return a.Lvalue
}

func (a *AssignVarStmt) GetRValue() expr.Rvalue {
	return a.Rvalue
}

func (a *AssignVarStmt) GetLineNumber() int {
	return a.LineNumber
}

func (a *AssignVarStmt) String() string {
	return a.Lvalue.String() + "=" + a.Rvalue.String()
}

func (a *AssignVarStmt) Accept(visitor StmtVisitor) {
	visitor.VisitAssignVarStmt(*a)
}

func (a *AssignVarStmt) GetStmtSeq() int {
	return a.Seq
}

func (a *AssignVarStmt) SetStmtSeq(i int) {
	a.Seq = i
}
