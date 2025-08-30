package parser

import (
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/ast"
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/conf"
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/version"
	"github.com/0xundef/php-webshell-detect/frontend/php/version/php5"
	"github.com/0xundef/php-webshell-detect/frontend/php/version/php7"
	"github.com/0xundef/php-webshell-detect/frontend/php/version/scanner"
	"errors"
)

var (
	// ErrVersionOutOfRange is returned if the version is not supported
	ErrVersionOutOfRange = errors.New("the version is out of supported range")

	php5RangeStart = &version.Version{Major: 5}
	php5RangeEnd   = &version.Version{Major: 5, Minor: 6}

	php7RangeStart = &version.Version{Major: 7}
	php7RangeEnd   = &version.Version{Major: 7, Minor: 4}
)

// Parser interface
type Parser interface {
	Parse() int
	GetRootNode() ast.Vertex
}

func Parse(src []byte, config conf.Config) (ast.Vertex, error) {
	var parser Parser

	if config.Version == nil {
		config.Version = php7RangeEnd
	}

	if config.Version.InRange(php5RangeStart, php5RangeEnd) {
		lexer := scanner.NewLexer(src, config)
		parser = php5.NewParser(lexer, config)
		parser.Parse()
		return parser.GetRootNode(), nil
	}

	if config.Version.InRange(php7RangeStart, php7RangeEnd) {
		lexer := scanner.NewLexer(src, config)
		parser = php7.NewParser(lexer, config)
		parser.Parse()
		return parser.GetRootNode(), nil
	}

	return nil, ErrVersionOutOfRange
}
