package context

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
)

type Pointer interface {
	GetPointsToSet() *PointsToSet
	SetPointsToSet(set *PointsToSet)
	GetObjects() []*CSObj
	GetPointerType() lang.Type
	AddOutEdge(edge PointerFlowEdge)
	GetOutEdges() []PointerFlowEdge
	GetKey() string
	String() string
	Diff(pts *PointsToSet) *PointsToSet

	SetMergeCount(i int)
	GetMergeCount() int
}
