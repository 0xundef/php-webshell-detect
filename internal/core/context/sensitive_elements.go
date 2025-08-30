package context

import (
	"github.com/0xundef/php-webshell-detect/internal/core/heap"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
)

type CSCallSite struct {
	CallSite  stmt.InvokeStmt
	Ctx       Context
	container *CSMethod
}

func (c *CSCallSite) String() string {
	return c.Ctx.String() + ":" + c.CallSite.String()
}

type CSMethod struct {
	Method *lang.PHPMethod
	Ctx    Context
	Edges  []CallEdge
}

func (m *CSMethod) AddEdge(edge CallEdge) {
	m.Edges = append(m.Edges, edge)
}

type CSObj struct {
	Obj     heap.Obj
	Context Context
}

func (c *CSObj) String() string {
	return c.Context.String() + c.Obj.GetKey()
}

func (c *CSObj) GetKey() string {
	return c.Context.String() + c.Obj.GetKey()
}
