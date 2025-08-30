package core

import (
	"github.com/0xundef/php-webshell-detect/internal/core/common/queue"
	"github.com/0xundef/php-webshell-detect/internal/core/context"
	"strings"
)

// WList receive New or literal assign
// $a = New A;
// $b = "c";
type WorkList struct {
	PointerEntries map[string]*context.PointsToSet //map[Pointer]PointsToSet
	PointerMap     map[string]context.Pointer
	CallEdges      queue.Queue
}

func Map(elements []*context.CSObj) []string {
	var ret = make([]string, len(elements))
	for i, element := range elements {
		ret[i] = element.String()
	}
	return ret
}
func (w *WorkList) Sprint() string {
	var ret string
	for point, pts := range w.PointerEntries {
		ret += "[" + point + "->" + strings.Join(Map(pts.Set0.ToSlice()), ",") + "],"
	}
	return strings.TrimSuffix(ret, ",")
}
func (w *WorkList) AddEntry(pointer context.Pointer, pts *context.PointsToSet) {
	if pts0, ok := w.PointerEntries[pointer.GetKey()]; ok {
		pts0.AddAll(pts)
	} else {
		w.PointerEntries[pointer.GetKey()] = &context.PointsToSet{Set0: pts.Set0.Clone()}
		w.PointerMap[pointer.GetKey()] = pointer
	}
}

func (w *WorkList) PollEntry() Entry {
	if !w.CallEdges.Empty() {
		if items, err := w.CallEdges.Poll(1, 0); err == nil {
			callEdge := items[0]
			return &CallEdgeEntry{
				CallEdge: callEdge.(context.CallEdge),
			}
		}
	} else if len(w.PointerEntries) > 0 {
		for key, set := range w.PointerEntries {
			entry := PointerEntry{
				Pointer:     w.PointerMap[key],
				PointsToSet: set,
			}
			delete(w.PointerEntries, key)
			delete(w.PointerMap, key)
			return &entry
		}
	}
	return nil
}
