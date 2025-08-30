package test

import (
	"github.com/0xundef/php-webshell-detect/internal/core"
	"fmt"
	"testing"
)

func Test_sample1(t *testing.T) {
	result := core.BatchDetect("../sample/webshell_data/", "php", "../config/")
	fmt.Println(result.String())
}

func Test_sample2(t *testing.T) {
	result := core.BatchDetect("../sample/WebShell_Black/", "php", "../config/")
	fmt.Println(result.String())
}
