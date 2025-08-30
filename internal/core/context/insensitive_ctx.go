package context

import (
	"github.com/0xundef/php-webshell-detect/internal/core/heap"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
)

type ContextInsensitiveSelector struct {
	ContextSelector
}

func (s *ContextInsensitiveSelector) GetEmptyContext() Context {
	return ContextEmpty()
}
func (s *ContextInsensitiveSelector) SelectContext0(callSite *CSCallSite, callee *lang.PHPMethod) Context {
	return s.GetEmptyContext()
}
func (s *ContextInsensitiveSelector) SelectContext1(callSite *CSCallSite, recv *CSObj, callee *lang.PHPMethod) Context {
	return s.GetEmptyContext()
}
func (s *ContextInsensitiveSelector) SelectHeapContext(method *CSMethod, obj heap.Obj) Context {
	return s.GetEmptyContext()
}
func (s *ContextInsensitiveSelector) SelectNewObjContext(method CSMethod, obj heap.Obj) Context {
	return s.GetEmptyContext()
}

func CreateCISelector() *ContextInsensitiveSelector {
	return &ContextInsensitiveSelector{}
}
