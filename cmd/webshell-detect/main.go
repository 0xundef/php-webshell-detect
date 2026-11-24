package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/0xundef/php-webshell-detect/internal/core"
	"github.com/0xundef/php-webshell-detect/internal/core/common"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	// Build-time variables (injected by the Makefile / GoReleaser).
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GoVersion = "unknown"
)

func init() {
	if GoVersion == "unknown" {
		GoVersion = runtime.Version()
	}
}

const (
	exitClean    = 0
	exitFindings = 1
	exitUsage    = 2
	exitError    = 3
)

type options struct {
	json      bool
	quiet     bool
	debug     bool
	version   bool
	configDir string
}

// jsonReport is the stable machine-readable output of one analysis.
type jsonReport struct {
	Path       string     `json:"path"`
	Tags       []string   `json:"tags"`
	Confidence string     `json:"confidence"`
	Findings   findings   `json:"findings"`
	LoopCount  int        `json:"loop_count"`
	CostMillis int64      `json:"cost_ms"`
	Version    string     `json:"version"`
	Error      string     `json:"error,omitempty"`
	Meta       reportMeta `json:"meta"`
}

type findings struct {
	SourceSink []string `json:"source_sink,omitempty"`
	RegEval    []string `json:"reg_eval,omitempty"`
	DynamicFun []string `json:"dynamic_fun,omitempty"`
	Decode     []string `json:"decode,omitempty"`
}

type reportMeta struct {
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
	GoVersion string `json:"go_version"`
	ConfDir   string `json:"config_dir"`
}

const usageText = `webshell-detect - static PHP webshell detection

Usage:
  webshell-detect [options] <file.php>

Options:
  -json            emit a JSON report on stdout
  -quiet           suppress informational logging
  -debug           enable verbose solver logging
  -config <dir>    directory containing conf.yaml and mock.php
  -version         print version information and exit
  -h, -help        show this help

Exit codes:
  0  analysis completed, no finding
  1  analysis completed, at least one finding
  2  usage or configuration error
  3  analysis error

Examples:
  webshell-detect /path/to/suspicious.php
  webshell-detect -json -config ./config /path/to/suspicious.php
`

func printUsage(w io.Writer) { fmt.Fprint(w, usageText) }

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	for _, arg := range args {
		if arg == "-h" || arg == "-help" || arg == "--help" {
			printUsage(os.Stdout)
			return exitClean
		}
	}

	opts := options{}
	fs := flag.NewFlagSet("webshell-detect", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.BoolVar(&opts.json, "json", false, "emit a JSON report on stdout")
	fs.BoolVar(&opts.quiet, "quiet", false, "suppress informational logging")
	fs.BoolVar(&opts.debug, "debug", false, "enable verbose solver logging")
	fs.BoolVar(&opts.version, "version", false, "print version information and exit")
	fs.StringVar(&opts.configDir, "config", "", "directory containing conf.yaml and mock.php")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, usageText)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return exitClean
		}
		return exitUsage
	}
	if opts.version {
		fmt.Printf("github.com/0xundef/php-webshell-detect %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Go Version: %s\n", GoVersion)
		return exitClean
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "error: exactly one PHP file must be supplied")
		fs.Usage()
		return exitUsage
	}

	configureLogging(opts)
	target := fs.Arg(0)

	confDir, err := resolveConfigDir(opts.configDir)
	if err != nil {
		return fail(opts, target, err, exitUsage)
	}
	conf, err := loadConfig(confDir)
	if err != nil {
		return fail(opts, target, err, exitUsage)
	}
	conf.ConfigDir = confDir
	conf.Path = target
	conf.Debug = (conf.Debug || opts.debug) && !opts.quiet

	start := time.Now()
	world, err := core.New(conf)
	if err != nil {
		return fail(opts, target, err, exitError)
	}
	item, err := world.AnalyzeSafe()
	item.CostTime = int64(time.Since(start).Milliseconds())
	if err != nil {
		return fail(opts, target, err, exitError)
	}

	if opts.json {
		writeJSON(item, confDir)
	} else {
		fmt.Printf("Analysis Result for %s:\n", target)
		fmt.Printf("Path: %s\n", item.Path)
		fmt.Printf("Confidence: %s\n", item.Confidence)
		fmt.Printf("Tags: %s\n", item.Tags())
		fmt.Printf("Version: %s\n", Version)
	}

	if item.HasAny() {
		return exitFindings
	}
	return exitClean
}

func configureLogging(opts options) {
	common.InitDefaultLogger()
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	switch {
	case opts.debug:
		logrus.SetLevel(logrus.DebugLevel)
	case opts.quiet:
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
	}
}

// resolveConfigDir finds the directory holding conf.yaml and mock.php. Order:
// explicit flag, environment variable, ./config, then <executable>/config.
func resolveConfigDir(explicit string) (string, error) {
	candidates := make([]string, 0, 4)
	if explicit != "" {
		candidates = append(candidates, explicit)
	}
	if env := os.Getenv("WEBSHELL_DETECT_CONFIG"); env != "" {
		candidates = append(candidates, env)
	}
	candidates = append(candidates, "config")
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "config"))
	}

	for _, dir := range candidates {
		if dir == "" {
			continue
		}
		if info, err := os.Stat(filepath.Join(dir, "conf.yaml")); err == nil && !info.IsDir() {
			if _, err := os.Stat(filepath.Join(dir, "mock.php")); err == nil {
				return dir, nil
			}
		}
	}
	if explicit != "" {
		return "", fmt.Errorf("config directory %q does not contain conf.yaml and mock.php", explicit)
	}
	return "", fmt.Errorf("cannot find conf.yaml/mock.php: use -config or run from the project root")
}

func loadConfig(dir string) (core.DetectConf, error) {
	conf := core.DetectConf{}
	data, err := os.ReadFile(filepath.Join(dir, "conf.yaml"))
	if err != nil {
		return conf, fmt.Errorf("read conf.yaml: %w", err)
	}
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return conf, fmt.Errorf("parse conf.yaml: %w", err)
	}
	return conf, nil
}

func writeJSON(item core.AnalyzeItem, confDir string) {
	report := jsonReport{
		Path:       item.Path,
		Tags:       item.TagList(),
		Confidence: item.Confidence,
		Findings: findings{
			SourceSink: item.SourceSink,
			RegEval:    item.RegEval,
			DynamicFun: item.DynamicFun,
			Decode:     item.Decode,
		},
		LoopCount:  item.LoopCount,
		CostMillis: item.CostTime,
		Version:    Version,
		Meta: reportMeta{
			BuildTime: BuildTime,
			GitCommit: GitCommit,
			GoVersion: GoVersion,
			ConfDir:   confDir,
		},
	}
	if report.Tags == nil {
		report.Tags = []string{}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(report); err != nil {
		fmt.Fprintf(os.Stderr, "error: encode report: %v\n", err)
	}
}

func fail(opts options, path string, err error, code int) int {
	if opts.json {
		report := jsonReport{
			Path:       path,
			Tags:       []string{},
			Confidence: core.ConfidenceNone,
			Version:    Version,
			Error:      err.Error(),
			Meta:       reportMeta{BuildTime: BuildTime, GitCommit: GitCommit, GoVersion: GoVersion},
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(report)
	}
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	return code
}
