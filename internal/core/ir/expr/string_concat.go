package expr

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

// StringConcatExp $a = "aa"."bb"
type StringConcatExp struct {
	BinaryExpr
	Operator BinaryOp
	Lop      Var
	Rop      Var
}

func (s *StringConcatExp) GetOperator() BinaryOp {
	return s.Operator
}

func (s *StringConcatExp) GetLOperand() Var {
	return s.Lop
}

func (s *StringConcatExp) GetROperand() Var {
	return s.Rop
}

func (s *StringConcatExp) getPrimitiveType() lang.Type {
	return s.Rop.Type
}

func (s *StringConcatExp) String() string {
	return s.Lop.String() + s.GetOperator().String() + s.Rop.String()
}
