/*

A Parser for PHP written in Go

Package usage example:

	package main

	import (
		"log"
		"os"

		"github.com/0xundef/php-webshell-detect/frontend/php/pkg/conf"
		"github.com/0xundef/php-webshell-detect/frontend/php/pkg/errors"
		"github.com/0xundef/php-webshell-detect/frontend/php/pkg/parser"
		"github.com/0xundef/php-webshell-detect/frontend/php/pkg/version"
		"github.com/0xundef/php-webshell-detect/frontend/php/pkg/visitor/dumper"
	)

	func main() {
		src := []byte(`<? echo "Hello world";`)

		// Error handler

		var parserErrors []*errors.Error
		errorHandler := func(e *errors.Error) {
			parsmakeerErrors = append(parserErrors, e)
		}

		// Parse

		rootNode, err := parser.Parse(src, conf.Config{
			Version:          &version.Version{Major: 5, Minor: 6},
			ErrorHandlerFunc: errorHandler,
		})

		if err != nil {
			log.Fatal("Error:" + err.Error())
		}

		// Dump

		goDumper := dumper.NewDumper(os.Stdout).
			WithTokens().
			WithPositions()

		rootNode.Accept(goDumper)
	}
*/
package parser
