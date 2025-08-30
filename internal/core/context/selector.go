package context

import (
	"github.com/0xundef/php-webshell-detect/internal/core/heap"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
)

type ContextSelector interface {
	GetEmptyContext() Context
	SelectContext0(callSite *CSCallSite, callee *lang.PHPMethod) Context
	SelectContext1(callSite *CSCallSite, recv *CSObj, callee *lang.PHPMethod) Context
	SelectHeapContext(method *CSMethod, obj heap.Obj) Context
	SelectNewObjContext(method CSMethod, obj heap.Obj) Context
}
