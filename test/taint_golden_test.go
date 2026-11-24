package test

import (
	"path/filepath"
	"testing"

	"github.com/0xundef/php-webshell-detect/internal/core"
)

// TestTaintGolden pins the detection output for the regression corpus. Each
// entry documents the expected tags and confidence for one fixture.
func TestTaintGolden(t *testing.T) {
	cases := []struct {
		file       string
		tags       string
		confidence string
	}{
		{"case1.php", "SourceSink", "high"},
		{"case2.php", "SourceSink", "high"},
		{"case4.php", "SourceSink,SuspiciousDecode", "high"},
		{"case5.php", "SourceSink", "high"},
		{"case6.php", "SourceSink", "high"},
		{"case7.php", "SourceSink", "high"},
		{"case8.php", "DynamicFun", "medium"},
		{"case9.php", "SourceSink", "high"},
		{"case10.php", "SourceSink", "high"},
		{"case11.php", "SourceSink", "high"},
		{"case12.php", "SourceSink", "high"},
		{"case13.php", "SourceSink", "high"},
		{"case14.php", "SourceSink", "high"},
		{"case15.php", "RegEval", "high"},
		{"case16.php", "SourceSink", "high"},
		{"case17.php", "SourceSink", "high"},
		{"case18.php", "SourceSink", "high"},
		{"case19.php", "SourceSink,SuspiciousDecode", "high"},
		{"case21.php", "DynamicFun", "medium"},
		{"case22.php", "SourceSink,SuspiciousDecode", "high"},
		{"case23.php", "DynamicFun", "medium"},
		{"case24.php", "DynamicFun", "medium"},
		{"case25.php", "DynamicFun", "medium"},
		{"case27.php", "DynamicFun", "medium"},
		{"case30.php", "DynamicFun", "medium"},
		{"case31.php", "DynamicFun", "medium"},
		// PHP function names are case-insensitive.
		{"todo4.php", "SourceSink", "high"},
	}

	base := core.DetectConf{
		ContextPolicy: "k-callsite-2",
		LoopLimit:     100000,
		ConfigDir:     "../config/",
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			conf := base
			conf.Path = filepath.Join("taint", tc.file)

			world, err := core.New(conf)
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			item, err := world.AnalyzeSafe()
			if err != nil {
				t.Fatalf("Analyze: %v", err)
			}
			if got := item.Tags(); got != tc.tags {
				t.Errorf("tags = %q, want %q (findings: %+v)", got, tc.tags, item)
			}
			if item.Confidence != tc.confidence {
				t.Errorf("confidence = %q, want %q", item.Confidence, tc.confidence)
			}
		})
	}
}

// TestKnownGaps documents fixtures the current engine does not detect yet.
// These are skipped rather than asserted so that fixing them is a test update,
// not a silent regression.
func TestKnownGaps(t *testing.T) {
	for _, file := range []string{"todo1.php", "todo2.php", "todo6.php"} {
		file := file
		t.Run(file, func(t *testing.T) {
			t.Skip("known limitation: reflection / __destruct / chr() obfuscation not modelled")
		})
	}
}

func TestBenignFiles(t *testing.T) {
	base := core.DetectConf{ContextPolicy: "k-callsite-2", LoopLimit: 100000, ConfigDir: "../config/"}
	for _, file := range []string{"benign_constant.php", "benign_function_arg.php"} {
		file := file
		t.Run(file, func(t *testing.T) {
			conf := base
			conf.Path = filepath.Join("taint", file)
			world, err := core.New(conf)
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			item, err := world.AnalyzeSafe()
			if err != nil {
				t.Fatalf("Analyze: %v", err)
			}
			if item.HasAny() {
				t.Errorf("benign file flagged: %+v", item)
			}
			if item.Confidence != core.ConfidenceNone {
				t.Errorf("confidence = %q, want %q", item.Confidence, core.ConfidenceNone)
			}
		})
	}
}
