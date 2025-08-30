package ir

type PHPContext struct {
	global map[string]bool
}

func (p *PHPContext) IsGlobal(name string) bool {
	if _, ok := p.global[name]; ok {
		return true
	}
	return false
}

var PHP_CONTEXT *PHPContext

func InitPHPContext() *PHPContext {
	return &PHPContext{
		global: map[string]bool{
			"$_GET":     false,
			"$_POST":    false,
			"$_COOKIE":  false,
			"$_REQUEST": false,
			"$_SERVER":  false,
		},
	}
}
