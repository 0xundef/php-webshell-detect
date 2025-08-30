package heap

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

type MockObj struct {
	Obj
	type0 lang.Type
	name  string
	Tags  string
}

func (o *MockObj) GetObjType() lang.Type {
	return o.type0
}
func (o *MockObj) GetKey() string {
	return o.name
}
func (o *MockObj) String() string {
	return o.name
}

func (o *MockObj) SetTag(tag string) {
	o.Tags = o.Tags + tag
}

func (o *MockObj) GetTags() string {
	return o.Tags
}
