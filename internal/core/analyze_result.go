package core

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/0xundef/php-webshell-detect/internal/core/common"
	"gopkg.in/yaml.v3"
)

// AnalyzeResult stores every file's result in batch analysis mode.
type AnalyzeResult struct {
	Found    AItems
	NotFound []string
	Error    []string
	CostTime int64
}

func (a AnalyzeResult) String() string {
	a.SortFounds()
	foundCount := len(a.Found)
	notFound := len(a.NotFound)
	errCount := len(a.Error)
	total := foundCount + notFound + errCount
	valid := foundCount + notFound

	rate := 0.0
	if valid > 0 {
		rate = float64(foundCount) / float64(valid)
	}

	ret := fmt.Sprintf("total: %v(%v)\nfound: %v\nnot found: %v\nerror: %v\nrate: %.2f\n", total, valid, foundCount, notFound, errCount, rate)
	ret += fmt.Sprintf("total cost time: %v ms\n", a.CostTime)

	ret += "top cost time file:\n"
	for _, item := range a.TopN(50) {
		ret += fmt.Sprintf("%v\tcost time(ms):%v loop count:%v\n", item.Path, item.CostTime, item.LoopCount)
	}

	var loops string
	for _, item := range a.Found {
		if item.LoopCount > 100000 {
			loops += fmt.Sprintf("%v\tcost time(ms):%v loop count:%v\n", item.Path, item.CostTime, item.LoopCount)
		}
	}
	if loops != "" {
		ret += "many loop count:\n" + loops
	}
	return ret
}

// MoveNotFound copies the files that were not flagged into test/notfound.
func (a AnalyzeResult) MoveNotFound() error {
	const dst = "test/notfound"
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, file := range a.NotFound {
		if err := common.Cp(file, filepath.Join(dst, filepath.Base(file))); err != nil {
			return fmt.Errorf("copy %s: %w", file, err)
		}
	}
	return nil
}

// SaveToFile writes the findings matching tags to path and returns how many
// findings were written. It does not print to stdout.
func (a AnalyzeResult) SaveToFile(path string, tags HIT_TAG) (count int, err error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return 0, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	writer := bufio.NewWriter(file)
	for _, item := range a.Found {
		if tags == All || item.HitTags&tags != 0 {
			count++
			if _, werr := writer.WriteString(item.String()); werr != nil {
				return count, werr
			}
		}
	}
	if _, err = writer.WriteString("done"); err != nil {
		return count, err
	}
	if err = writer.Flush(); err != nil {
		return count, err
	}
	return count, nil
}

func (a AnalyzeResult) SortFounds() {
	sort.Sort(a.Found)
}

// TopN returns up to n findings ordered by cost, or all findings when n is not
// smaller than the result size.
func (a AnalyzeResult) TopN(n int) []AnalyzeItem {
	if n <= 0 || n >= len(a.Found) {
		return a.Found
	}
	return a.Found[:n]
}

// BatchDetect analyses every file below dir with the given suffix. It is meant
// for test/corpus runs, not for the production CLI. A missing directory or a
// bad configuration is reported in Result.Error instead of panicking.
func BatchDetect(dir string, suffix string, configDir string, excludes ...string) AnalyzeResult {
	result := AnalyzeResult{
		Found:    []AnalyzeItem{},
		NotFound: []string{},
		Error:    []string{},
	}
	start := time.Now()

	conf := DetectConf{}
	confData, err := os.ReadFile(filepath.Join(configDir, "conf.yaml"))
	if err != nil {
		result.Error = append(result.Error, fmt.Sprintf("read conf.yaml: %v", err))
		return result
	}
	if err = yaml.Unmarshal(confData, &conf); err != nil {
		result.Error = append(result.Error, fmt.Sprintf("parse conf.yaml: %v", err))
		return result
	}
	conf.ConfigDir = configDir

	if _, statErr := os.Stat(dir); statErr != nil {
		result.Error = append(result.Error, fmt.Sprintf("scan directory %q: %v", dir, statErr))
		result.CostTime = time.Since(start).Milliseconds()
		return result
	}

	walkErr := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			result.Error = append(result.Error, fmt.Sprintf("%s: %v", path, err))
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, suffix) {
			return nil
		}
		for _, exclude := range excludes {
			if exclude != "" && strings.Contains(path, exclude) {
				return nil
			}
		}

		fileStart := time.Now()
		cfg := conf
		cfg.Path = path
		world, err := New(cfg)
		if err != nil {
			result.Error = append(result.Error, fmt.Sprintf("%s: %v", path, err))
			return nil
		}
		item, err := world.AnalyzeSafe()
		item.Path = path
		item.CostTime = time.Since(fileStart).Milliseconds()
		if err != nil {
			result.Error = append(result.Error, fmt.Sprintf("%s: %v", path, err))
			return nil
		}
		if item.HasAny() {
			result.Found = append(result.Found, item)
		} else {
			result.NotFound = append(result.NotFound, path)
		}
		return nil
	})
	if walkErr != nil {
		result.Error = append(result.Error, fmt.Sprintf("walk %q: %v", dir, walkErr))
	}

	result.CostTime = time.Since(start).Milliseconds()
	return result
}
