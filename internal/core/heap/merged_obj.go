package heap

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"strconv"
	"strings"
)

type MergedObj struct {
	Obj
	Id    int
	Tags  string
	type0 lang.Type
}

func (o *MergedObj) GetObjType() lang.Type {
	return o.type0
}

func (o *MergedObj) GetKey() string {
	return "merged" + strconv.Itoa(o.Id) + "_" + o.Tags
}
func (o *MergedObj) String() string {
	return "merged" + strconv.Itoa(o.Id) + "_" + o.Tags
}

func (o *MergedObj) GetTags() string {
	return o.Tags
}
func (o *MergedObj) SetTag(tag string) {
	if !strings.Contains(o.Tags, tag) {
		o.Tags = o.Tags + tag
	}
}
