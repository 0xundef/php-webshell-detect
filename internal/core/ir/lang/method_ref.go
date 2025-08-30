package lang

type MethodRef struct {
	MemberRef
	DeclaringClass *PHPClass
	Name           string

	IsDynamic  bool
	ReturnType Type
	Signature  string
}

func (ref *MethodRef) ResolveMethod() *PHPMethod {
	//use HA to find
	return nil
}

func (ref *MethodRef) String() string {
	if ref.IsDynamic {

	}
	return ref.Name
}
func (ref *MethodRef) GetReturnType() Type {
	if ref.IsDynamic {

	}
	return ref.ReturnType
}
