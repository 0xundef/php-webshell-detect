package php7_test

import (
	"io/ioutil"
	"testing"

	"github.com/0xundef/php-webshell-detect/frontend/php/version/php7"
	"github.com/0xundef/php-webshell-detect/frontend/php/version/scanner"
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/conf"
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/version"
)

func BenchmarkPhp7(b *testing.B) {
	src, err := ioutil.ReadFile("test.php")

	if err != nil {
		b.Fatal("can not read test.php: " + err.Error())
	}

	for n := 0; n < b.N; n++ {
		config := conf.Config{
			Version: &version.Version{
				Major: 7,
				Minor: 4,
			},
		}
		lexer := scanner.NewLexer(src, config)
		php7parser := php7.NewParser(lexer, config)
		php7parser.Parse()
	}
}
