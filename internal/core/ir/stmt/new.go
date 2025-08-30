package stmt

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
)

// NewStmt $shell = new COM("WScript.shell")
type NewStmt struct {
	DefinitionStmt
	Lvalue     expr.Var
	Rvalue     expr.NewExpr
	LineNumber int
	Seq        int
}

func (n *NewStmt) GetLValue() expr.Lvalue {
	return &n.Lvalue
}

func (n *NewStmt) GetRValue() expr.Rvalue {
	return n.Rvalue
}

func (n *NewStmt) GetLineNumber() int {
	return n.LineNumber
}

func (n *NewStmt) String() string {
	return n.Lvalue.String() + "=" + n.Rvalue.String()
}

func (n *NewStmt) Accept(visitor StmtVisitor) {
	visitor.VisitNewStmt(*n)
}

func (n *NewStmt) GetStmtSeq() int {
	return n.Seq
}

func (n *NewStmt) SetStmtSeq(i int) {
	n.Seq = i
}
