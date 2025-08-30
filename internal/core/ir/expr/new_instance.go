package expr

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

type NewExpr interface {
	Rvalue
	IsNewExpr() bool
}

type NewInstanceExpr struct {
	NewExpr
	Type0 lang.ClassType
}

func (n *NewInstanceExpr) String() string {
	return "new " + n.Type0.Name
}
func (n *NewInstanceExpr) GetType() lang.Type {
	return &n.Type0
}
