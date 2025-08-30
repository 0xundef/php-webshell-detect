package stmt

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
)

// AssignBinaryStmt $a='a'^'b'
type AssignBinaryStmt struct {
	DefinitionStmt
	Lvalue     expr.Lvalue
	Rvalue     expr.Rvalue
	LineNumber int
	Seq        int
}

func (a *AssignBinaryStmt) GetLValue() expr.Lvalue {
	return a.Lvalue
}

func (a *AssignBinaryStmt) GetRValue() expr.Rvalue {
	return a.Rvalue
}

func (a *AssignBinaryStmt) GetLineNumber() int {
	return a.LineNumber
}

func (a *AssignBinaryStmt) String() string {
	return a.Lvalue.String() + "=" + a.Rvalue.String()
}

func (a *AssignBinaryStmt) Accept(visitor StmtVisitor) {
	visitor.VisitAssignAssignBinaryStmt(*a)
}

func (a *AssignBinaryStmt) GetStmtSeq() int {
	return a.Seq
}

func (a *AssignBinaryStmt) SetStmtSeq(i int) {
	a.Seq = i
}
