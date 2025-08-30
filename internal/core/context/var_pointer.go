package context

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
	"strconv"
)

// CSVarPointer subtype of Pointer
type CSVarPointer struct {
	Context     Context
	PointsToSet *PointsToSet
	successors  []Pointer
	outEdges    []PointerFlowEdge

	var0        expr.Var
	invokes     []stmt.InvokeStmt
	loadFields  []stmt.LoadFieldStmt
	storeFields []stmt.StoreFieldStmt
	loadArrays  []StmtContext
	mergeObject []StmtContext //binary
	storeArrays []stmt.StoreArrayStmt

	IsFunName  bool
	MergeCount int
}

type StmtContext struct {
	Stmt stmt.Stmt
	Ctx  Context
}

func (c *CSVarPointer) SetMergeCount(i int) {
	c.MergeCount = i
}
func (c *CSVarPointer) GetMergeCount() int {
	return c.MergeCount
}

func (c *CSVarPointer) IsSuspiciousDynamicFun() string {
	if c.IsFunName && c.MergeCount > 5 {
		return strconv.Itoa(c.MergeCount)
	}
	return ""
}

func (c *CSVarPointer) GetVar() expr.Var {
	return c.var0
}
func (c *CSVarPointer) GetInvokeStmts() []stmt.InvokeStmt {
	return c.invokes
}
func (c *CSVarPointer) AddInvokeStmt(invoke stmt.InvokeStmt) {
	c.invokes = append(c.invokes, invoke)
}

func (c *CSVarPointer) GetLoadFields() []stmt.LoadFieldStmt {
	return c.loadFields
}
func (c *CSVarPointer) AddLoadField(load stmt.LoadFieldStmt) {
	c.loadFields = append(c.loadFields, load)
}
func (c *CSVarPointer) AddStoreField(store stmt.StoreFieldStmt) {
	c.storeFields = append(c.storeFields, store)
}
func (c *CSVarPointer) GetStoreField() []stmt.StoreFieldStmt {
	return c.storeFields
}

func (c *CSVarPointer) AddStoreArray(store stmt.StoreArrayStmt) {
	c.storeArrays = append(c.storeArrays, store)
}
func (c *CSVarPointer) GetStoreArrays() []stmt.StoreArrayStmt {
	return c.storeArrays
}

func (c *CSVarPointer) GetMerge() []StmtContext {
	return c.mergeObject
}

func (c *CSVarPointer) AddMerge(stmt stmt.AssignBinaryStmt, context Context) {
	item := StmtContext{
		Stmt: &stmt,
		Ctx:  context,
	}
	c.mergeObject = append(c.mergeObject, item)
}

func (c *CSVarPointer) AddLoadArray(stmt stmt.LoadArrayStmt, context Context) {
	item := StmtContext{
		Stmt: &stmt,
		Ctx:  context,
	}
	c.loadArrays = append(c.loadArrays, item)
}
func (c *CSVarPointer) GetLoadArrays() []StmtContext {
	return c.loadArrays
}

func (c *CSVarPointer) GetPointsToSet() *PointsToSet {
	return c.PointsToSet
}
func (c *CSVarPointer) SetPointsToSet(set *PointsToSet) {
	c.PointsToSet = set
}
func (c *CSVarPointer) GetObjects() []*CSObj {
	return c.PointsToSet.GetObjects()
}
func (c *CSVarPointer) GetPointerType() lang.Type {
	return c.var0.GetType()
}
func (c *CSVarPointer) AddOutEdge(edge PointerFlowEdge) {
	c.outEdges = append(c.outEdges, edge)
}
func (c *CSVarPointer) GetOutEdges() []PointerFlowEdge {
	return c.outEdges
}
func (c *CSVarPointer) GetKey() string {
	return c.Context.String() + c.var0.Name
}
func (c *CSVarPointer) String() string {
	return c.Context.String() + c.var0.Name
}
func (c *CSVarPointer) Diff(pts *PointsToSet) *PointsToSet {
	if c.GetPointsToSet() == nil {
		c.PointsToSet = pts
		return pts
	}
	return c.GetPointsToSet().AddDiff(pts)
}
