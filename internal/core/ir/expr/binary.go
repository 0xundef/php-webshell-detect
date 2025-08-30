package expr

type BinaryExpr interface {
	Rvalue
	GetOperator() BinaryOp
	GetLOperand() Var
	GetROperand() Var
}

type BinaryOp int

const (
	OR BinaryOp = iota
	AND
	XOR
	CONCAT
)

func (op BinaryOp) String() string {
	return [...]string{"|", "&", "^", "."}[op]
}
