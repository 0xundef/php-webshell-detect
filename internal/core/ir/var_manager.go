package ir

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
)

type VarManager struct {
	Consts      map[string]interface{} // todo remove
	GlobalVars  map[string]interface{}
	GlobalArray map[string]interface{}
}

func (m *VarManager) AddConst(stmt *stmt.AssignLiteralStmt) {
	if _, ok := m.Consts[stmt.Lvalue.String()]; !ok {
		m.Consts[stmt.Lvalue.String()] = stmt.Rvalue.String()
	}
}

func (m *VarManager) GetConst(v expr.Var) interface{} {
	return m.Consts[v.Name]
}

func (m *VarManager) GetGlobalVar(name string) interface{} {
	if v, ok := m.GlobalVars[name]; ok {
		return v
	}
	return nil
}

func (m *VarManager) AddGlobalVar(name string, v expr.Var) {
	m.GlobalVars[name] = v
}

func CreateVarManager() *VarManager {
	return &VarManager{
		Consts:      map[string]interface{}{},
		GlobalVars:  map[string]interface{}{},
		GlobalArray: map[string]interface{}{},
	}
}
