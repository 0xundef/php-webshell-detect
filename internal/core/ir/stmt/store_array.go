package stmt

import "github.com/0xundef/php-webshell-detect/internal/core/ir/expr"

// StoreArrayStmt $a["cmd"] = "xxx"
type StoreArrayStmt struct {
	DefinitionStmt
	Lvalue     expr.ArrayAccess
	Rvalue     expr.Var
	LineNumber int
	Seq        int
}

func (l *StoreArrayStmt) GetLineNumber() int {
	return l.LineNumber
}

func (l *StoreArrayStmt) String() string {
	return l.Lvalue.String() + "=" + l.Rvalue.String()
}

func (l *StoreArrayStmt) Accept(visitor StmtVisitor) {
	visitor.VisitStoreArrayStmt(*l)
}

func (l *StoreArrayStmt) GetLValue() expr.Lvalue {
	return &l.Lvalue
}

func (l *StoreArrayStmt) GetRValue() expr.Lvalue {
	return l.Rvalue
}
func (l *StoreArrayStmt) GetStmtSeq() int {
	return l.Seq
}

func (l *StoreArrayStmt) SetStmtSeq(i int) {
	l.Seq = i
}
