package core

type DetectConf struct {
	Path          string
	Debug         bool `yaml:"debug"`
	ConfigDir     string
	ContextPolicy string `yaml:"context_policy"`
	LoopLimit     int    `yaml:"loop_limit"`
	DataType      int    `yaml:"data_type"`
}
