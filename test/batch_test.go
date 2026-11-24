package test

import (
	"os"
	"testing"

	"github.com/0xundef/php-webshell-detect/internal/core"
)

func Test_sample1(t *testing.T) {
	const dir = "../sample/webshell_data/"
	if _, err := os.Stat(dir); err != nil {
		t.Skipf("sample corpus not available: %v", err)
	}
	result := core.BatchDetect(dir, "php", "../config/")
	t.Log(result.String())
}

func Test_sample2(t *testing.T) {
	const dir = "../sample/WebShell_Black/"
	if _, err := os.Stat(dir); err != nil {
		t.Skipf("sample corpus not available: %v", err)
	}
	result := core.BatchDetect(dir, "php", "../config/")
	t.Log(result.String())
}

func TestBatchDetectMissingDir(t *testing.T) {
	result := core.BatchDetect("../sample/does-not-exist", "php", "../config/")
	if len(result.Error) == 0 {
		t.Fatal("expected an error for a missing directory")
	}
}

func TestBatchDetectBadConfig(t *testing.T) {
	result := core.BatchDetect("../test/taint", "php", "../does-not-exist")
	if len(result.Error) == 0 {
		t.Fatal("expected an error for a missing config directory")
	}
}
