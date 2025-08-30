package lang

type PHPClassField struct {
	DeclaringClass *PHPClass
	Name           string

	Type0 Type
}

func (f *PHPClassField) GetFieldType() Type {
	return f.Type0
}

func (f *PHPClassField) String() string {
	return ""
}
