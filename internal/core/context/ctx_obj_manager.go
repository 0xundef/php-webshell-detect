package context

import "github.com/0xundef/php-webshell-detect/internal/core/heap"

type ObjContextKey struct {
	obj     string
	context string
}
type CSObjManager struct {
	objs map[ObjContextKey]*CSObj
}

func (m *CSObjManager) GetCSObj(ctx Context, obj heap.Obj) *CSObj {
	key := ObjContextKey{
		obj:     obj.GetKey(),
		context: ctx.String(),
	}
	if o, ok := m.objs[key]; ok {
		return o
	} else {
		m.objs[key] = &CSObj{
			Obj:     obj,
			Context: ctx,
		}
		return m.objs[key]
	}
}
