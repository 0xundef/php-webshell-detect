package core

import "github.com/0xundef/php-webshell-detect/internal/core/context"

// Entry the WorkList element, only two types:PointerEntry CallEdgeEntry
type Entry interface {
	IsEntry() bool
}

type PointerEntry struct {
	Entry
	Pointer     context.Pointer
	PointsToSet *context.PointsToSet
}

func (p *PointerEntry) IsEntry() bool {
	return true
}
func (p *PointerEntry) GetPointer() context.Pointer {
	return p.Pointer
}
func (p *PointerEntry) GetPointsToSet() *context.PointsToSet {
	return p.PointsToSet
}

type CallEdgeEntry struct {
	Entry
	CallEdge context.CallEdge
}
