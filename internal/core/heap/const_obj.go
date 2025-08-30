package heap

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"strconv"
)

// ConstObj /* subtype of Pointer
type ConstObj struct {
	Value expr.Literal
	Id    int
	Tags  string
	Count int
}

func (o *ConstObj) GetObjType() lang.Type {
	return o.Value.GetType()
}

func (o *ConstObj) GetKey() string {
	return "const" + strconv.Itoa(o.Id) + "_" + o.Value.String()
}
func (o *ConstObj) String() string {
	return o.Value.String()
}
func (o *ConstObj) GetTags() string {
	return o.Tags
}
func (o *ConstObj) SetTag(tag string) {
	o.Tags = o.Tags + tag
}
