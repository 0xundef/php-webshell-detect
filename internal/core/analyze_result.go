package core

import (
	"bufio"
	"github.com/0xundef/php-webshell-detect/internal/core/common"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// AnalyzeResult store all file's analyzed result in batch analyze mode
type AnalyzeResult struct {
	Found    AItems
	NotFound []string
	Error    []string
	CostTime int64
}

func (a AnalyzeResult) String() string {
	var ret string
	a.SortFounds()
	foundCount := len(a.Found)
	notFound := len(a.NotFound)
	err := len(a.Error)
	total := foundCount + notFound + err
	valid := foundCount + notFound
	rate := float64(foundCount) / float64(foundCount+notFound)

	count := fmt.Sprintf("total: %v(%v)\nfound: %v\nnot found: %v \nerror: %v \nrate: %.2f\n", total, valid, foundCount, notFound, err, rate)
	time := fmt.Sprintf("total cost time: %v sec", a.CostTime)

	var tops string
	topn := a.TopN(50)
	for _, item := range topn {
		tops = tops + fmt.Sprintf("%v	cost time:%v loop count:%v \n", item.Path, item.CostTime, item.LoopCount)
	}

	var loops string
	for _, item := range a.Found {
		if item.LoopCount > 100000 {
			loops = loops + fmt.Sprintf("%v	cost time:%v loop count:%v \n", item.Path, item.CostTime, item.LoopCount)
		}
	}
	fmt.Sprintf("top cost time file: \n")
	ret = count + "\n" + time + "\n" + tops + "\n" + " many loop count \n" + loops
	return ret
}

func (a AnalyzeResult) MoveNotFound() {
	common.Rmv("test/notfound")
	common.Mkdir("test/notfound")
	for _, file := range a.NotFound {
		common.Cp(file, "test/notfound")
	}
}
func (a AnalyzeResult) SaveToFile(path string, tags HIT_TAG) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// create a writer to write to the file
	writer := bufio.NewWriter(file)

	// write lines to the file
	var count int
	for _, item := range a.Found {
		if tags == All {
			writer.WriteString(item.String())
		} else if item.HitTags == tags {
			count++
			writer.WriteString(item.String())
		}
		if err != nil {
			panic(err)
		}
	}

	writer.WriteString("done")
	fmt.Printf(" save taints %v", count)

	// flush the writer to ensure all data is written to the file
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}

func (a AnalyzeResult) SortFounds() {
	sort.Sort(a.Found)
}

func (a AnalyzeResult) TopN(n int) []AnalyzeItem {
	if n > 0 && n < a.Found.Len() {
		return a.Found[:n]
	}
	return nil
}

// BatchDetect recommended for use this in test environments only
func BatchDetect(dir string, suffix string, configDir string, excludes ...string) AnalyzeResult {
	result := AnalyzeResult{
		Found:    []AnalyzeItem{},
		NotFound: []string{},
		Error:    []string{},
	}
	var count int64
	st := time.Now().Unix()
	countVsitor := func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			count++
		}
		return nil
	}
	filepath.Walk(dir, countVsitor)
	var bar *progressbar.ProgressBar

	wg := sync.WaitGroup{}
	conf := DetectConf{}
	dat, _ := ioutil.ReadFile(configDir + "/conf.yaml")
	err := yaml.Unmarshal(dat, &conf)
	conf.ConfigDir = configDir
	if err != nil {
		logrus.Error("read conf.yaml failed")
	}
	if conf.Debug {
		bar = progressbar.Default(count)
	}
	visitor := func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			for _, exclude := range excludes {
				if strings.Contains(path, exclude) {
					return nil
				}
			}
			wg.Add(1)

			go func() {
				defer func() {
					if conf.Debug {
						bar.Add(1)
					}
					wg.Done()

					if err := recover(); err != nil {
						result.Error = append(result.Error, path)
					}
				}()
				config := conf
				config.Path = path
				st := time.Now().Unix()
				world := Bang(config)
				taints := world.Analyze()
				ed := time.Now().Unix()

				if taints.HasAny() {
					taints.Path = path
					taints.CostTime = ed - st
					result.Found = append(result.Found, taints)
				} else {
					result.NotFound = append(result.NotFound, path)
				}
			}()
		}
		return nil
	}

	filepath.Walk(dir, visitor)

	wg.Wait()
	ed := time.Now().Unix()
	result.CostTime = ed - st
	return result
}
