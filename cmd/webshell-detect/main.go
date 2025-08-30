package main

import (
	"github.com/0xundef/php-webshell-detect/internal/core"
	"github.com/0xundef/php-webshell-detect/internal/core/common"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var (
	// Build-time variables (injected by Makefile)
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GoVersion = "unknown"
	
	Conf *core.DetectConf
)

func init() {
	runtime.GOMAXPROCS(4)
	rand.Seed(time.Now().UnixNano())
	Conf = &core.DetectConf{}
	dat, _ := ioutil.ReadFile("./config/conf.yaml")
	err := yaml.Unmarshal(dat, Conf)
	Conf.ConfigDir = "./config/"
	if err != nil {
		common.Log().Error("./config/conf.yaml load error")
	}
}

func check(path string, conf *core.DetectConf) (ret core.AnalyzeItem, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				common.Log().Error(err)
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	conf.Path = path
	world := core.Bang(*conf)
	ret = world.Analyze()
	return ret, nil
}



func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github.com/0xundef/php-webshell-detect <file_path>")
		fmt.Println("       github.com/0xundef/php-webshell-detect --version")
		fmt.Println("       github.com/0xundef/php-webshell-detect --help")
		os.Exit(1)
	}

	// Handle version flag
	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Printf("github.com/0xundef/php-webshell-detect %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Go Version: %s\n", GoVersion)
		os.Exit(0)
	}

	// Handle help flag
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		fmt.Println("github.com/0xundef/php-webshell-detect - A tool for detecting PHP webshells using static analysis")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  github.com/0xundef/php-webshell-detect <file_path>    Analyze a PHP file for webshell patterns")
		fmt.Println("  github.com/0xundef/php-webshell-detect --version      Show version information")
		fmt.Println("  github.com/0xundef/php-webshell-detect --help         Show this help message")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  github.com/0xundef/php-webshell-detect /path/to/suspicious.php")
		fmt.Println("  github.com/0xundef/php-webshell-detect --version")
		os.Exit(0)
	}

	filePath := os.Args[1]
	common.Log().Info("php webShell detect started...")
	
	ret, err := check(filePath, Conf)
	if err != nil {
		common.Log().Errorf("Error analyzing file %s: %v", filePath, err)
		os.Exit(1)
	}

	fmt.Printf("Analysis Result for %s:\n", filePath)
	fmt.Printf("Path: %s\n", ret.Path)
	fmt.Printf("Confidence: %s\n", ret.Confidence)
	fmt.Printf("Tags: %s\n", ret.Tags())
	fmt.Printf("Version: %s\n", Version)
	
	common.Log().Info("php webShell detect completed.")
}
