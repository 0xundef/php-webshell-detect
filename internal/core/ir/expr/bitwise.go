package expr

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

type BitwiseExp struct {
	BinaryExpr
	Operator BinaryOp
	Lop      Var
	Rop      Var
}

func (b *BitwiseExp) GetOperator() BinaryOp {
	return b.Operator
}

func (b *BitwiseExp) GetLOperand() Var {
	return b.Lop
}

func (b *BitwiseExp) GetROperand() Var {
	return b.Rop
}

func (b *BitwiseExp) GetType() lang.Type {
	return b.Rop.Type
}

func (b *BitwiseExp) String() string {
	return b.Lop.String() + b.Operator.String() + b.Rop.String()
}
