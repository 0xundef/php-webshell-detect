package test

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// item in b do not exist in a
func sub(a []string, b []string) []string {
	var result []string
	setA := make(map[string]bool)
	for _, v := range a {
		setA[v] = false
	}
	for _, v := range b {
		setA[v] = true
	}
	for k, v := range setA {
		if !v {
			result = append(result, k)
		}
	}
	return result
}

func load(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Initialize an empty array to store the lines
	var lines []string

	// Iterate over each line of the file and append it to the array
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Check if there was an error while scanning the file
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return nil
	}

	// Print the contents of the array
	return lines
}

func extract(dir string, suffix string) []string {
	var ret []string
	visitor := func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			ret = append(ret, path)
		}
		return nil
	}
	filepath.Walk(dir, visitor)
	return ret
}
