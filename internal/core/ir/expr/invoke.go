package expr

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
)

type Invoke interface {
	Rvalue
	GetContainer() *lang.PHPMethod
	GetArgs() []Var
	GetFunName() string
}
