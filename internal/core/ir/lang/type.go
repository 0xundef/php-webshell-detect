package lang

type Type interface {
	String() string
}

type ReferenceType interface {
	Type
}

type ScalarType interface {
	Type
}

type UnKnowType struct {
	Type
}

func (u *UnKnowType) String() string {
	return "UnKnowType"
}
