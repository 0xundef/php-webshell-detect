package conf

import (
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/errors"
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/version"
)

type Config struct {
	Version          *version.Version
	ErrorHandlerFunc func(e *errors.Error)
}
