package context

// CallEdge subtype of Edge
type CallEdge struct {
	CallSite *CSCallSite
	Callee   *CSMethod
}

func (c *CallEdge) String() string {
	return c.CallSite.String() + "**" + c.CallSite.String()
}
