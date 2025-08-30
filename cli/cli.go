package main

import (
	"github.com/0xundef/php-webshell-detect/internal/core"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var path, configDir string

	flag.StringVar(&path, "path", "", "the file path to detect")
	flag.StringVar(&configDir, "config_dir", "", "the conf file path")
	flag.Parse()

	var usage = `
usage of webshell_detect:
--path="":the file path to detect
--config_dir="":the conf file path`

	if path == "" {
		fmt.Println(usage)
		return
	}
	conf := core.DetectConf{}
	dat, err := ioutil.ReadFile(configDir + "/conf.yaml")
	if err != nil {
		fmt.Println("read conf.yaml failed")
		return
	}
	yaml.Unmarshal(dat, &conf)
	conf.Path = path
	conf.ConfigDir = configDir
	world := core.Bang(conf)
	taints := world.Analyze()
	fmt.Println(taints.String())
}
