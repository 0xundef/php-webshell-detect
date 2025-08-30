package stmt

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
)

// Stmt supper interface for all statements
type Stmt interface {
	GetLineNumber() int
	GetStmtSeq() int
	SetStmtSeq(i int)
	String() string
	Accept(visitor StmtVisitor)
}
type DefinitionStmt interface {
	Stmt
	GetLValue() expr.Lvalue
	GetRValue() expr.Rvalue
}
