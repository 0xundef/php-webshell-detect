package common

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func TrimString(value string) string {
	if value[0:1] == "\"" {
		return strings.Trim(value, "\"")
	} else if value[0:1] == "'" {
		return strings.Trim(value, "'")
	}
	return value
}
func Cp(file string, path string) {
	cmd := exec.Command("cp", file, path)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("执行命令失败：", err)
		return
	}
	fmt.Println(string(output))
}
func Rmv(path string) {
	cmd := exec.Command("rm", "-rf", path)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("执行命令失败：", err)
		return
	}
	fmt.Println(string(output))
}
func Mkdir(path string) {
	cmd := exec.Command("mkdir", "-p", path)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("执行命令失败：", err)
		return
	}
	fmt.Println(string(output))
}

func ExtractInputs(html string) []string {
	var ret []string
	re := regexp.MustCompile(`<input\s+[^>]*name="([^"]+)"[^>]*>`)
	match := re.FindAllStringSubmatch(html, 5)

	for _, input := range match {
		if strings.Contains(input[0], `type="text"`) {
			ret = append(ret, input[1])
		}
	}
	return ret
}
