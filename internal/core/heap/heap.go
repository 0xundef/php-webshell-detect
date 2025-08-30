package heap

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
	"strconv"
)

type HeapModel interface {
	GetConstObj(value expr.Literal) Obj
	GetNewObj(value stmt.NewStmt) Obj
	GetMockObj(name string) Obj
	GetMergedObj(stmt stmt.Stmt) *MergedObj
}

type DefaultHeapModel struct {
	HeapModel
	ConstObjs  map[string]*ConstObj
	NewObjs    map[string]*NewObj
	MergedObjs map[string]*MergedObj
	Count      int
}

func CreateHeapModel() HeapModel {
	return &DefaultHeapModel{
		ConstObjs: map[string]*ConstObj{},
		NewObjs:   map[string]*NewObj{},
		MergedObjs: map[string]*MergedObj{
			"string": &MergedObj{
				Id:    0,
				Tags:  "",
				type0: nil,
			},
		},
	}
}

func (h *DefaultHeapModel) GetMergedObj(stmt stmt.Stmt) *MergedObj {
	return h.MergedObjs["string"]
}

func (h *DefaultHeapModel) GetConstObj(value expr.Literal) Obj {
	t := value.GetType()
	key := t.String() + value.String()
	if obj, ok := h.ConstObjs[key]; ok {
		return obj
	}
	h.Count++
	obj := &ConstObj{
		Value: value,
		Id:    h.Count,
	}
	h.ConstObjs[key] = obj
	return obj
}

func (h *DefaultHeapModel) GetNewObj(newSite stmt.NewStmt) Obj {
	key := strconv.Itoa(newSite.GetLineNumber()) + newSite.Rvalue.GetType().String()
	if obj, ok := h.NewObjs[key]; ok {
		return obj
	}
	h.Count++
	newObj := &NewObj{
		Value: newSite.Rvalue,
		Id:    h.Count,
	}
	h.NewObjs[key] = newObj
	return newObj
}

func (h *DefaultHeapModel) GetMockObj(name string) Obj {
	mockObj := &MockObj{name: name}
	return mockObj
}
