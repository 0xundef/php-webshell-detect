package test

import (
	"github.com/0xundef/php-webshell-detect/internal/core"
	"fmt"
	"testing"
)

var config = core.DetectConf{
	Debug:         false,
	ContextPolicy: "k-callsite-2",
	LoopLimit:     100000,
	ConfigDir:     "../config/",
}

//ob_start callback
func Test_case1(t *testing.T) {
	config.Path = "./taint/case1.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case2(t *testing.T) {
	config.Path = "./taint/case2.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case3(t *testing.T) {
	config.Path = "./taint/case3.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case4(t *testing.T) {
	config.Path = "./taint/case4.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case5(t *testing.T) {
	config.Path = "./taint/case5.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case6(t *testing.T) {
	config.Path = "./taint/case6.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case7(t *testing.T) {
	config.Path = "./taint/case7.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case8(t *testing.T) {
	config.Path = "./taint/case8.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case9(t *testing.T) {
	config.Path = "./taint/case9.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case10(t *testing.T) {
	config.Path = "./taint/case10.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case11(t *testing.T) {
	config.Path = "./taint/case11.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case12(t *testing.T) {
	config.Path = "./taint/case12.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case13(t *testing.T) {
	config.Path = "./taint/case13.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case14(t *testing.T) {
	config.Path = "./taint/case14.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}
func Test_case15(t *testing.T) {
	config.Path = "./taint/case15.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case16(t *testing.T) {
	config.Path = "./taint/case16.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case17(t *testing.T) {
	config.Path = "./taint/case17.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case18(t *testing.T) {
	config.Path = "./taint/case18.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case19(t *testing.T) {
	config.Path = "./taint/case19.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case20(t *testing.T) {
	config.Path = "./taint/case20.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case21(t *testing.T) {
	config.Path = "./taint/case21.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case22(t *testing.T) {
	config.Path = "./taint/case22.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case23(t *testing.T) {
	config.Path = "./taint/case23.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case24(t *testing.T) {
	config.Path = "./taint/case24.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}

func Test_case25(t *testing.T) {
	config.Path = "./taint/case25.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints)
}

func Test_case26(t *testing.T) {
	config.Path = "./taint/case26.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints)
}

func Test_case27(t *testing.T) {
	config.Path = "./taint/case27.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints)
}

func Test_case28(t *testing.T) {
	config.Path = "./taint/case28.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints)
}

func Test_case29(t *testing.T) {
	config.Path = "./taint/case29.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints)
}

func Test_case30(t *testing.T) {
	config.Path = "./taint/case30.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints)
}

func Test_case31(t *testing.T) {
	config.Path = "./taint/case31.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints)
}

func Test_case32(t *testing.T) {
	config.Path = "./taint/case32.php"
	world := core.Bang(config)
	taints := world.Analyze()
	fmt.Println(taints.String())
}
