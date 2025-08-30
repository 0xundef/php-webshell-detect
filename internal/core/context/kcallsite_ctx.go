package context

import (
	"github.com/0xundef/php-webshell-detect/internal/core/heap"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
)

type KcallSiteSelector struct {
	ContextSelector
	limit int
}

func (s *KcallSiteSelector) GetEmptyContext() Context {
	return ContextEmpty()
}
func (s *KcallSiteSelector) SelectContext0(callSite *CSCallSite, callee *lang.PHPMethod) Context {
	return ContextAppend(callSite.Ctx, callSite.CallSite, s.limit)
}
func (s *KcallSiteSelector) SelectContext1(callSite *CSCallSite, recv *CSObj, callee *lang.PHPMethod) Context {
	return ContextAppend(callSite.Ctx, callSite.CallSite, s.limit)
}
func (s *KcallSiteSelector) SelectHeapContext(method *CSMethod, obj heap.Obj) Context {
	return s.GetEmptyContext()
}
func (s *KcallSiteSelector) SelectNewObjContext(method CSMethod, obj heap.Obj) Context {
	return s.GetEmptyContext()
}

func CreateKCallSiteSelector(limit int) *KcallSiteSelector {
	return &KcallSiteSelector{
		limit: limit,
	}
}
