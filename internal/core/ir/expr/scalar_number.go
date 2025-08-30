package expr

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"strconv"
)

type ScalarNumber struct {
	Literal
	Value int
}

func (n *ScalarNumber) String() string {
	return strconv.Itoa(n.Value)
}
func (n *ScalarNumber) GetType() lang.Type {
	return &lang.ScalarNumberType{}
}
