package lang

type ScalarNumberType struct {
	ScalarType
}

func (t *ScalarNumberType) String() string {
	return "number"
}
