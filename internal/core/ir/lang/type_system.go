package lang

type TypeSystem interface {
	GetType(name string) Type
	GetClassType(name string) ClassType
	IsSubType(super Type, sub Type) bool
}

type PHPTypeSystem struct {
	TypeSystem
	Hierarchy *PHPClassHierarchy
}

func (ts *PHPTypeSystem) GetType(name string) Type {
	return nil
}
func (ts *PHPTypeSystem) GetClassType(name string) ClassType {
	return ClassType{
		ReferenceType: nil,
		Name:          "",
		Class:         &PHPClass{},
	}
}
func (ts *PHPTypeSystem) IsSubType(super Type, sub Type) bool {
	return true
}
