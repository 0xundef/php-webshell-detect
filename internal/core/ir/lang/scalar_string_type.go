package lang

type ScalarStringType struct {
	ScalarType
}

func (t *ScalarStringType) String() string {
	return "string"
}
