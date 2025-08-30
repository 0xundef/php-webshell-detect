package context

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

type CSMethodManager struct {
	methodMap map[string]*CSMethod
	methods   []*CSMethod
}

func (m *CSMethodManager) GetKey(ctx Context, method *lang.PHPMethod) string {
	return ctx.String() + method.GetSignature()
}

func (m *CSMethodManager) GetCSMethod(ctx Context, method *lang.PHPMethod) *CSMethod {
	key := m.GetKey(ctx, method)
	if cm, ok := m.methodMap[key]; ok {
		return cm
	}

	csMethod := &CSMethod{
		Method: method,
		Ctx:    ctx,
		Edges:  []CallEdge{},
	}
	m.methodMap[key] = csMethod
	m.methods = append(m.methods, csMethod)
	return m.methodMap[key]
}
