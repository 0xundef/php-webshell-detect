package context

import (
	"github.com/0xundef/php-webshell-detect/internal/core/heap"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
)

/*
CSManager is the final view that filled with all the var -> obj mappings
it is usually called in diff
*/
type CSManager struct {
	PrtManager    *PointerManager
	ObjManager    *CSObjManager
	MethodManager *CSMethodManager
}

func (m *CSManager) GetCSVar(context Context, var0 expr.Var) *CSVarPointer {
	return m.PrtManager.GetCSVar(context, var0)
}

func (m *CSManager) GetInstanceField(base *CSObj, field *lang.PHPClassField) *InstanceFieldPointer {
	return m.PrtManager.GetInstanceField(*base, *field)
}

func (m *CSManager) GetArrayPointer(array *CSObj) *ArrayPointer {
	return m.PrtManager.GetArrayPointer(*array)
}

func (m *CSManager) GetCSObj(context Context, obj heap.Obj) *CSObj {
	return m.ObjManager.GetCSObj(context, obj)
}

func (m *CSManager) GetCSCallSite(ctx Context, invoke stmt.InvokeStmt) *CSCallSite {
	csContainer := m.MethodManager.GetCSMethod(ctx, invoke.Container)
	return &CSCallSite{
		CallSite:  invoke,
		Ctx:       ctx,
		container: csContainer,
	}
}

func (m *CSManager) GetCSMethod(ctx Context, method *lang.PHPMethod) *CSMethod {
	return m.MethodManager.GetCSMethod(ctx, method)
}

func CreateCSManager() *CSManager {
	methodManager := CSMethodManager{
		methodMap: map[string]*CSMethod{},
		methods:   make([]*CSMethod, 1),
	}
	ptrManager := PointerManager{
		Vars:           map[string]*CSVarPointer{},
		InstanceFields: map[string]*InstanceFieldPointer{},
		Arrays:         map[string]*ArrayPointer{},
	}
	objManager := CSObjManager{objs: map[ObjContextKey]*CSObj{}}
	return &CSManager{
		PrtManager:    &ptrManager,
		ObjManager:    &objManager,
		MethodManager: &methodManager,
	}
}
