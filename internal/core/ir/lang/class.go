package lang

import "strings"

type PHPClass struct {
	Name       string
	Namespace  string
	Type0      *ClassType
	Super      *PHPClass
	Interfaces []PHPClass

	DeclaredFields  map[string]*PHPClassField
	DeclaredMethods map[string]*PHPMethod
}

// AddMethod registers a method. PHP method names are case-insensitive, so the
// lookup key is normalised while the method keeps its original spelling.
func (class *PHPClass) AddMethod(method *PHPMethod) {
	class.DeclaredMethods[strings.ToLower(method.Name)] = method
}

func (class *PHPClass) AddField(field *PHPClassField) {
	class.DeclaredFields[field.Name] = field
}

func (class *PHPClass) GetDeclaredMethod(name string) *PHPMethod {
	return class.DeclaredMethods[strings.ToLower(name)]
}

// GetDeclaredField returns the field, creating a placeholder entry when it has
// not been declared explicitly (implicit public property).
func (class *PHPClass) GetDeclaredField(name string) *PHPClassField {
	field := class.DeclaredFields[name]
	if field == nil {
		class.DeclaredFields[name] = &PHPClassField{
			DeclaringClass: class,
			Name:           name,
			Type0:          nil,
		}
		return class.DeclaredFields[name]
	}
	return field
}

type ClassMember interface {
}

func CreatePHPClass(name string, namespace string) *PHPClass {
	return &PHPClass{
		Name:            name,
		Namespace:       namespace,
		Type0:           nil,
		Super:           nil,
		Interfaces:      nil,
		DeclaredFields:  map[string]*PHPClassField{},
		DeclaredMethods: map[string]*PHPMethod{},
	}
}
