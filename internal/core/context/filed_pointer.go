package context

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

// InstanceFieldPointer subtype of Pointer
type InstanceFieldPointer struct {
	base        CSObj
	field       lang.PHPClassField
	PointsToSet *PointsToSet
	successors  []Pointer
	outEdges    []PointerFlowEdge
	MergeCount  int
}

func (f *InstanceFieldPointer) SetMergeCount(i int) {
	f.MergeCount = i
}
func (f *InstanceFieldPointer) GetMergeCount() int {
	return f.MergeCount
}
func (f *InstanceFieldPointer) GetPointsToSet() *PointsToSet {
	return f.PointsToSet
}
func (f *InstanceFieldPointer) SetPointsToSet(set *PointsToSet) {
	f.PointsToSet = set
}
func (f *InstanceFieldPointer) GetObjects() []*CSObj {
	return f.PointsToSet.Set0.ToSlice()
}
func (f *InstanceFieldPointer) GetPointerType() lang.Type {
	return f.field.Type0
}
func (f *InstanceFieldPointer) AddOutEdge(edge PointerFlowEdge) {
	f.outEdges = append(f.outEdges, edge)
}
func (f *InstanceFieldPointer) GetOutEdges() []PointerFlowEdge {
	return f.outEdges
}
func (f *InstanceFieldPointer) GetKey() string {
	return f.base.String() + "->" + f.field.Name
}
func (f *InstanceFieldPointer) String() string {
	return f.base.String() + "->" + f.field.Name
}
func (f *InstanceFieldPointer) Diff(pts *PointsToSet) *PointsToSet {
	if f.GetPointsToSet() == nil {
		f.PointsToSet = pts
		return pts
	}
	return f.GetPointsToSet().AddDiff(pts)
}
