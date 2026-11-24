package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// captureStdout runs fn while capturing everything written to os.Stdout.
func captureStdout(t *testing.T, fn func() int) (string, int) {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	code := fn()
	w.Close()
	os.Stdout = old
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(data), code
}

func configDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.Abs("../../config")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeTarget(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "target.php")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunJSONFinding(t *testing.T) {
	target := writeTarget(t, "<?php\nsystem($_GET['c']);\n")
	out, code := captureStdout(t, func() int {
		return run([]string{"-json", "-quiet", "-config", configDir(t), target})
	})
	if code != exitFindings {
		t.Fatalf("exit = %d, want %d (output: %s)", code, exitFindings, out)
	}
	var report jsonReport
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if report.Path != target {
		t.Errorf("path = %q, want %q", report.Path, target)
	}
	if len(report.Tags) != 1 || report.Tags[0] != "SourceSink" {
		t.Errorf("tags = %v, want [SourceSink]", report.Tags)
	}
	if report.Confidence != "high" {
		t.Errorf("confidence = %q, want high", report.Confidence)
	}
}

func TestRunCleanFile(t *testing.T) {
	target := writeTarget(t, "<?php\nsystem(\"ls -l\");\n")
	out, code := captureStdout(t, func() int {
		return run([]string{"-json", "-quiet", "-config", configDir(t), target})
	})
	if code != exitClean {
		t.Fatalf("exit = %d, want %d (output: %s)", code, exitClean, out)
	}
	var report jsonReport
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if len(report.Tags) != 0 {
		t.Errorf("tags = %v, want empty", report.Tags)
	}
	if report.Confidence != "none" {
		t.Errorf("confidence = %q, want none", report.Confidence)
	}
}

func TestRunMissingConfig(t *testing.T) {
	target := writeTarget(t, "<?php\n")
	_, code := captureStdout(t, func() int {
		return run([]string{"-json", "-quiet", "-config", filepath.Join(t.TempDir(), "nope"), target})
	})
	if code != exitUsage {
		t.Fatalf("exit = %d, want %d", code, exitUsage)
	}
}

func TestRunMissingTarget(t *testing.T) {
	_, code := captureStdout(t, func() int {
		return run([]string{"-json", "-quiet", "-config", configDir(t), filepath.Join(t.TempDir(), "missing.php")})
	})
	if code != exitError {
		t.Fatalf("exit = %d, want %d", code, exitError)
	}
}

func TestRunUsage(t *testing.T) {
	if code := run([]string{}); code != exitUsage {
		t.Fatalf("exit = %d, want %d", code, exitUsage)
	}
	if code := run([]string{"-version"}); code != exitClean {
		t.Fatalf("version exit = %d, want %d", code, exitClean)
	}
}
