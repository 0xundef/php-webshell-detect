package lang

type ArrayType struct {
	ReferenceType
	Name string
}

func (t *ArrayType) String() string {
	return t.Name
}
func (t *ArrayType) IsReferenceType() bool {
	return true
}
