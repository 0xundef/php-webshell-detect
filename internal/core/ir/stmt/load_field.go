package stmt

import "github.com/0xundef/php-webshell-detect/internal/core/ir/expr"

// LoadFieldStmt $obj->age = "xxx"
type LoadFieldStmt struct {
	DefinitionStmt
	Lvalue     expr.Var
	Rvalue     expr.FieldAccess
	LineNumber int
	Seq        int
}

func (l *LoadFieldStmt) GetLineNumber() int {
	return l.LineNumber
}

func (l *LoadFieldStmt) String() string {
	return l.Lvalue.String() + "=" + l.Rvalue.String()
}

func (l *LoadFieldStmt) Accept(visitor StmtVisitor) {
	visitor.VisitLoadFieldStmt(*l)
}

func (l *LoadFieldStmt) GetLValue() expr.Lvalue {
	return l.Lvalue
}

func (l *LoadFieldStmt) GetRValue() expr.Lvalue {
	return &l.Rvalue
}

func (l *LoadFieldStmt) GetStmtSeq() int {
	return l.Seq
}

func (l *LoadFieldStmt) SetStmtSeq(i int) {
	l.Seq = i
}
