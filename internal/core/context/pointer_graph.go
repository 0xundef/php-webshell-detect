package context

type PointerFlowEdge struct {
	Edge
	Source Pointer
	Target Pointer
	Kind   EdgeKind
}

func (p *PointerFlowEdge) GetSource() interface{} {
	return p.Source
}
func (p *PointerFlowEdge) GetTarget() interface{} {
	return p.Target
}
func (p *PointerFlowEdge) String() string {
	return p.Source.String() + "->" + p.Target.String()
}

type PointerGraph struct {
	Pointers map[string]Pointer
}

func (g *PointerGraph) hasNode(node Pointer) bool {
	return true
}
func (g *PointerGraph) hasEdge(node Pointer) bool {
	return false
}
func (g *PointerGraph) getPredsOf(node Pointer) []Pointer {
	//no need
	return nil
}
func (g *PointerGraph) getInEdgesOf(node Pointer) []Edge {
	//no need
	return nil
}
func (g *PointerGraph) getOutEdgesOf(node Pointer) []Edge {
	return nil
}

// specified
func (g *PointerGraph) getPointers() map[string]Pointer {
	return g.Pointers
}
func (g *PointerGraph) AddEdge(edge PointerFlowEdge) {
	fromKey := edge.Source.GetKey()
	toKey := edge.Target.GetKey()
	if _, ok := g.Pointers[fromKey]; !ok {
		g.Pointers[fromKey] = edge.Source
	}
	if _, ok := g.Pointers[toKey]; !ok {
		g.Pointers[toKey] = edge.Target
	}
	source := g.Pointers[fromKey]
	source.AddOutEdge(edge)
}
