package heap

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"strconv"
)

type NewObj struct {
	Obj
	Value expr.NewExpr
	Id    int
	Tags  string
}

func (o *NewObj) GetObjType() lang.Type {
	return o.Value.GetType()
}

func (o *NewObj) GetKey() string {
	return "newObj_" + o.Value.GetType().String() + "_" + strconv.Itoa(o.Id)
}
func (o *NewObj) String() string {
	return o.Value.String()
}

func (o *NewObj) SetTag(tag string) {
	o.Tags = o.Tags + tag
}

func (o *NewObj) GetTags() string {
	return o.Tags
}
