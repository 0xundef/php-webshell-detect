package common

import (
	"io"
	"os"
	"regexp"
	"strings"
)

// TrimString removes a surrounding pair of single or double quotes. It is safe
// for empty input (the previous implementation indexed value[0:1] and panicked).
func TrimString(value string) string {
	if len(value) < 2 {
		return value
	}
	first, last := value[0], value[len(value)-1]
	if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
		return value[1 : len(value)-1]
	}
	return value
}

// Cp copies a regular file to path.
func Cp(file string, path string) error {
	src, err := os.Open(file)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err = io.Copy(dst, src); err != nil {
		dst.Close()
		return err
	}
	return dst.Close()
}

// Rmv removes path recursively; missing paths are not an error.
func Rmv(path string) error {
	return os.RemoveAll(path)
}

// Mkdir creates path and its parents.
func Mkdir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// ExtractInputs extracts name attributes of text inputs from an HTML snippet.
func ExtractInputs(html string) []string {
	var ret []string
	re := regexp.MustCompile(`<input\s+[^>]*name="([^"]+)"[^>]*>`)
	for _, match := range re.FindAllStringSubmatch(html, 5) {
		if strings.Contains(match[0], `type="text"`) {
			ret = append(ret, match[1])
		}
	}
	return ret
}
