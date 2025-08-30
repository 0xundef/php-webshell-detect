package lang

type ClassHierarchy interface {
	ResolveMethod(mref MethodRef) *PHPMethod
	AllClasses() []PHPClass
}

type PHPClassHierarchy struct {
	ClassHierarchy
	Classes []*PHPClass
}

func (h *PHPClassHierarchy) ResolveMethod(ref MethodRef) *PHPMethod {
	ret := h.LookupMethod(ref.DeclaringClass, ref.Name)
	return ret
}

func (h *PHPClassHierarchy) LookupMethod(class *PHPClass, name string) *PHPMethod {
	var tmp *PHPClass
	tmp = class
	for true {
		if tmp != nil {
			m := tmp.GetDeclaredMethod(name)
			if m != nil {
				return m
			}
			tmp = tmp.Supper
		} else {
			break
		}
	}
	return nil
}

func (h *PHPClassHierarchy) LookupField(class *PHPClass, name string) *PHPClassField {
	var tmp *PHPClass
	tmp = class
	for true {
		if tmp != nil {
			m := tmp.GetDeclaredFiled(name)
			if m != nil {
				return m
			}
			tmp = tmp.Supper
		} else {
			break
		}
	}
	return nil
}

func (h *PHPClassHierarchy) AddClass(class *PHPClass) {
	h.Classes = append(h.Classes, class)
}
func (h *PHPClassHierarchy) GetClass(name string) *PHPClass {
	for _, class := range h.Classes {
		if name == class.Name {
			return class
		}
	}
	return nil
}

func CreateClassHierarchy() *PHPClassHierarchy {
	return &PHPClassHierarchy{
		Classes: []*PHPClass{},
	}
}
