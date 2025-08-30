package stmt

import "github.com/0xundef/php-webshell-detect/internal/core/ir/expr"

// LoadArrayStmt $a = $_GET["cmd"]
type LoadArrayStmt struct {
	DefinitionStmt
	Lvalue     expr.Var
	Rvalue     expr.ArrayAccess
	LineNumber int
	Seq        int
}

func (l *LoadArrayStmt) GetLineNumber() int {
	return l.LineNumber
}

func (l *LoadArrayStmt) String() string {
	return l.Lvalue.String() + "=" + l.Rvalue.String()
}

func (l *LoadArrayStmt) Accept(visitor StmtVisitor) {
	visitor.VisitLoadArrayStmt(*l)
}

func (l *LoadArrayStmt) GetLValue() expr.Lvalue {
	return l.Lvalue
}

func (l *LoadArrayStmt) GetRValue() expr.Rvalue {
	return &l.Rvalue
}

func (l *LoadArrayStmt) GetStmtSeq() int {
	return l.Seq
}

func (l *LoadArrayStmt) SetStmtSeq(i int) {
	l.Seq = i
}
