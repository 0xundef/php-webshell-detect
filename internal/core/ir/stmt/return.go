package stmt

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
)

// ReturnStmt return statement if exist
type ReturnStmt struct {
	Stmt
	Value      expr.Var
	LineNumber int
	Seq        int
}

func (r *ReturnStmt) GetLineNumber() int {
	return r.LineNumber
}

func (r *ReturnStmt) String() string {
	return "return " + r.Value.String()
}

func (r *ReturnStmt) Accept(visitor StmtVisitor) {
	//do nothing
}
func (r *ReturnStmt) GetValue() expr.Var {
	return r.Value
}
func (r *ReturnStmt) GetStmtSeq() int {
	return r.Seq
}
func (r *ReturnStmt) SetStmtSeq(i int) {
	r.Seq = i
}
