package expr

import "github.com/0xundef/php-webshell-detect/internal/core/ir/lang"

type ScalarString struct {
	Literal
	Value string
}

func (s *ScalarString) String() string {
	return s.Value
}
func (s *ScalarString) GetType() lang.Type {
	return &lang.ScalarStringType{}
}
