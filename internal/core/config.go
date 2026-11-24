package core

// DetectConf configures one analysis run.
type DetectConf struct {
	// Path is the PHP file to analyse.
	Path string
	// Debug enables verbose solver/IR logging.
	Debug bool
	// ConfigDir is the directory containing conf.yaml and mock.php.
	ConfigDir string
	// ContextPolicy selects the context sensitivity: "ci" or "k-callsite-<n>".
	ContextPolicy string `yaml:"context_policy"`
	// LoopLimit caps solver iterations; <= 0 means DefaultLoopLimit.
	LoopLimit int `yaml:"loop_limit"`
}
