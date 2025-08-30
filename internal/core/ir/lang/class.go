package lang

type PHPClass struct {
	Name       string
	Namespace  string
	Type0      *ClassType
	Supper     *PHPClass
	Interfaces []PHPClass

	DeclaredFields  map[string]*PHPClassField
	DeclaredMethods map[string]*PHPMethod
}

func (class *PHPClass) AddMethod(method *PHPMethod) {
	class.DeclaredMethods[method.Name] = method
}
func (class *PHPClass) AddField(field *PHPClassField) {
	class.DeclaredFields[field.Name] = field
}
func (class *PHPClass) GetDeclaredMethod(name string) *PHPMethod {
	return class.DeclaredMethods[name]
}
func (class *PHPClass) GetDeclaredFiled(name string) *PHPClassField {
	field := class.DeclaredFields[name]
	if field == nil {
		class.DeclaredFields[name] = &PHPClassField{
			DeclaringClass: class,
			Name:           name,
			Type0:          nil,
		}
		return class.DeclaredFields[name]
	} else {
		return field
	}
}

type ClassMember interface {
}

func CreatePHPClass(name string, namespace string) *PHPClass {
	return &PHPClass{
		Name:            name,
		Namespace:       namespace,
		Type0:           nil,
		Supper:          nil,
		Interfaces:      nil,
		DeclaredFields:  map[string]*PHPClassField{},
		DeclaredMethods: map[string]*PHPMethod{},
	}

}
