package context

import (
	"github.com/0xundef/php-webshell-detect/internal/core/heap"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
)

type KObjSelector struct {
	ContextSelector
}

func (s *KObjSelector) SelectContext0(callSite CSCallSite, callee lang.PHPMethod) Context {
	return callSite.Ctx
}
func (s *KObjSelector) SelectContext1(callSite CSCallSite, recv CSObj, callee lang.PHPMethod) Context {
	return ContextAppend(recv.Context, recv.Obj, 2)
}
func (s *KObjSelector) SelectHeapContext(method CSMethod, obj heap.Obj) Context {
	return s.GetEmptyContext()
}
func (s *KObjSelector) SelectNewObjContext(method CSMethod, obj heap.Obj) Context {
	return s.GetEmptyContext()
}
