package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTrimString(t *testing.T) {
	cases := map[string]string{
		"":      "",
		"abc":   "abc",
		`"abc"`: "abc",
		`'abc'`: "abc",
		`"`:     `"`,
		`""`:    "",
		`"a'b"`: "a'b",
		`'a"b'`: "a\"b"[:0] + "a\"b",
	}
	for in, want := range cases {
		if got := TrimString(in); got != want {
			t.Errorf("TrimString(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFileHelpers(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	dstDir := filepath.Join(dir, "nested", "dst")
	if err := Mkdir(dstDir); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dstDir, "copy.txt")
	if err := Cp(src, dst); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("copied content = %q", data)
	}
	if err := Rmv(dstDir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dstDir); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be removed", dstDir)
	}
}

func TestExtractInputs(t *testing.T) {
	html := `<form><input type="text" name="cmd"><input type="hidden" name="token"></form>`
	got := ExtractInputs(html)
	if len(got) != 1 || got[0] != "cmd" {
		t.Fatalf("ExtractInputs = %v, want [cmd]", got)
	}
}
