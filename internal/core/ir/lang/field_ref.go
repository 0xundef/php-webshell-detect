package lang

type FieldRef struct {
	MemberRef
	DeclaringClass PHPClass
	Name           string

	Field     PHPClassField
	Type0     Type
	IsStatic0 bool
}

func (f *FieldRef) IsStatic() bool {
	return f.IsStatic0
}
