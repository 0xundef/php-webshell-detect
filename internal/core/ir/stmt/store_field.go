package stmt

import "github.com/0xundef/php-webshell-detect/internal/core/ir/expr"

// StoreFieldStmt $obj->age = "xxx"
type StoreFieldStmt struct {
	DefinitionStmt
	Lvalue     expr.FieldAccess
	Rvalue     expr.Var
	LineNumber int
	Seq        int
}

func (l *StoreFieldStmt) GetLineNumber() int {
	return l.LineNumber
}

func (l *StoreFieldStmt) String() string {
	return l.Lvalue.String() + "=" + l.Rvalue.String()
}

func (l *StoreFieldStmt) Accept(visitor StmtVisitor) {
	visitor.VisitStoreFieldStmt(*l)
}

func (l *StoreFieldStmt) GetLValue() expr.Lvalue {
	return &l.Lvalue
}

func (l *StoreFieldStmt) GetRValue() expr.Lvalue {
	return l.Rvalue
}

func (l *StoreFieldStmt) GetStmtSeq() int {
	return l.Seq
}
func (l *StoreFieldStmt) SetStmtSeq(i int) {
	l.Seq = i
}
