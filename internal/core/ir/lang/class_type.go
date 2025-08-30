package lang

type ClassType struct {
	ReferenceType
	Name  string
	Class *PHPClass
}

func (t *ClassType) String() string {
	return t.Name
}
func (t *ClassType) IsReferenceType() bool {
	return true
}

func (t *ClassType) GetClass() PHPClass {
	return *t.Class
}
