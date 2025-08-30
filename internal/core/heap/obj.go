package heap

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

type Obj interface {
	GetObjType() lang.Type
	GetKey() string
	String() string
	SetTag(tag string)
	GetTags() string
}
