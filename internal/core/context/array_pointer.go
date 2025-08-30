package context

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

// ArrayPointer subtype of Pointer
type ArrayPointer struct {
	Array       CSObj
	PointsToSet *PointsToSet
	successors  []Pointer
	outEdges    []PointerFlowEdge

	MergeCount int
}

func (p *ArrayPointer) SetMergeCount(i int) {
	p.MergeCount = i
}
func (p *ArrayPointer) GetMergeCount() int {
	return p.MergeCount
}

func (p *ArrayPointer) GetPointsToSet() *PointsToSet {
	return p.PointsToSet
}
func (p *ArrayPointer) SetPointsToSet(set *PointsToSet) {
	p.PointsToSet = set
}
func (p *ArrayPointer) GetObjects() []*CSObj {
	return p.PointsToSet.Set0.ToSlice()
}
func (p *ArrayPointer) GetPointerType() lang.Type {
	return p.Array.Obj.GetObjType()
}
func (p *ArrayPointer) AddOutEdge(edge PointerFlowEdge) {
	p.outEdges = append(p.outEdges, edge)
}
func (p *ArrayPointer) GetOutEdges() []PointerFlowEdge {
	return p.outEdges
}
func (p *ArrayPointer) GetKey() string {
	return p.Array.String()
}
func (p *ArrayPointer) String() string {
	return p.Array.String() + "[*]"
}
func (p *ArrayPointer) Diff(pts *PointsToSet) *PointsToSet {
	if p.GetPointsToSet() == nil {
		p.PointsToSet = pts
		return pts
	}
	return p.GetPointsToSet().AddDiff(pts)
}
