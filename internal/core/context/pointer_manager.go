package context

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
	"github.com/deckarep/golang-set/v2"
)

type PointerManager struct {
	Vars           map[string]*CSVarPointer
	InstanceFields map[string]*InstanceFieldPointer
	Arrays         map[string]*ArrayPointer
}

func (m *PointerManager) GetCSVar(context Context, var0 expr.Var) *CSVarPointer {
	var key string
	if var0.Global {
		key = context.GetRoot().String() + var0.Name
	} else {
		key = context.String() + var0.Name
	}
	if v, ok := m.Vars[key]; ok {
		return v
	} else {
		m.Vars[key] = &CSVarPointer{
			Context:     context,
			var0:        var0,
			outEdges:    []PointerFlowEdge{},
			PointsToSet: &PointsToSet{Set0: mapset.NewSet[*CSObj]()},
			invokes:     []stmt.InvokeStmt{},
		}
		return m.Vars[key]
	}
}

func (m *PointerManager) GetInstanceField(base CSObj, filed lang.PHPClassField) *InstanceFieldPointer {
	key := base.String() + filed.Name
	if v, ok := m.InstanceFields[key]; ok {
		return v
	} else {
		m.InstanceFields[key] = &InstanceFieldPointer{
			base:        base,
			field:       filed,
			PointsToSet: &PointsToSet{Set0: mapset.NewSet[*CSObj]()},
		}
		return m.InstanceFields[key]
	}
}
func (m *PointerManager) GetArrayPointer(array CSObj) *ArrayPointer {
	key := array.String()
	if v, ok := m.Arrays[key]; ok {
		return v
	} else {
		m.Arrays[key] = &ArrayPointer{
			Array:       array,
			PointsToSet: &PointsToSet{Set0: mapset.NewSet[*CSObj]()},
		}
		return m.Arrays[key]
	}
}
