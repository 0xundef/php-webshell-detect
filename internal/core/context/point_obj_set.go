package context

import (
	"github.com/deckarep/golang-set/v2"
	"strings"
)

type PointsToSet struct {
	Set0 mapset.Set[*CSObj]
}

func (p *PointsToSet) AddObj(obj *CSObj) {
	p.Set0.Add(obj)
}

func (p *PointsToSet) AddAll(set *PointsToSet) {
	for _, obj := range set.Set0.ToSlice() {
		p.Set0.Add(obj)
	}
}

func (p *PointsToSet) AddDiff(set *PointsToSet) *PointsToSet {
	diff := set.Set0.Difference(p.Set0)
	for _, obj := range diff.ToSlice() {
		p.Set0.Add(obj)
	}
	return &PointsToSet{
		Set0: diff.Clone(),
	}
}

func (p *PointsToSet) Sprint() string {
	var ret []string
	for _, obj := range p.Set0.ToSlice() {
		ret = append(ret, obj.String())
	}
	return strings.Join(ret, ",")
}

func (p *PointsToSet) SprintIfHasEndWith(str string) string {
	var ret []string
	var hit bool
	for _, obj := range p.Set0.ToSlice() {
		if strings.HasSuffix(obj.String(), str) {
			hit = true
		}
		ret = append(ret, obj.String())
	}
	if hit {
		return strings.Join(ret, ",")
	} else {
		return ""
	}
}

func (p *PointsToSet) SprintIfHasEqual(str string) string {
	var ret []string
	var hit bool
	for _, obj := range p.Set0.ToSlice() {
		if strings.HasSuffix(obj.String(), "_"+str) {
			hit = true
		}
		ret = append(ret, obj.String())
	}
	if hit {
		return strings.Join(ret, ",")
	} else {
		return ""
	}
}

func (p *PointsToSet) IsEmpty() bool {
	return p.Set0.Cardinality() == 0
}
func (p *PointsToSet) GetObjects() []*CSObj {
	return p.Set0.ToSlice()
}
