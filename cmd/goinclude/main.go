package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	appName = "goinclude"
	version = "1.0"
)

var re = regexp.MustCompile(`{{\s*include\s+"([^"]+)"\s*}}`)

func main() {
	switch {
	case len(os.Args) == 1:
		help()
	case len(os.Args) == 3:
		includeFiles(os.Args[1], os.Args[2])
	default:
		log.Fatal(appName + ": provide input file path and output file path as parameters")
	}
}

func includeFiles(pathToTemplate, pathToOutputFile string) {
	tmpl, err := getFileContent(pathToTemplate)
	if err != nil {
		log.Fatal(appName + ": could not read " + pathToTemplate)
	}
	var includedFiles []string
	var get = func(filePath string) (s string, err error) {
		b, err := ioutil.ReadFile(getFullPath(filePath))
		if err != nil {
			return
		}
		includedFiles = append(includedFiles, filePath)
		return string(b), nil
	}
	out := re.ReplaceAllStringFunc(tmpl, func(s string) string {
		c, err := get(re.ReplaceAllString(s, "$1"))
		if err != nil {
			return fmt.Sprintf(
				"(%s: error including %s into %s: %s)",
				appName, getFullPath(pathToTemplate), pathToTemplate, err,
			)
		}
		return c
	})
	if len(includedFiles) > 0 {
		fmt.Println(strings.Join(includedFiles, " "))
	}
	if err := ioutil.WriteFile(pathToOutputFile, []byte(out), 0644); err != nil {
		log.Fatal(appName + ": error writing " + pathToOutputFile)
	}
}

func getFileContent(filePath string) (string, error) {
	b, err := ioutil.ReadFile(getFullPath(filePath))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getFullPath(pathToFile string) string {
	if !strings.HasPrefix(pathToFile, "/") {
		workingDirectory, err := os.Getwd()
		if err != nil {
			log.Fatal(fmt.Sprintf("%s: getting working directory: %s\n", appName, err))
		}
		pathToFile = path.Join(workingDirectory, pathToFile)
	}
	s, err := filepath.Abs(pathToFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s: resolving full path: %s\n", appName, err))
	}
	return s
}

func help() {
	fmt.Printf(`%s %s // github.com/zengabor/goinclude
Includes the content of the indicated files into a go template and outputs the merged result to the specified file.

This is how you include a file inside a template: {{ include "path/otherfile.html" }}

Usage:    %[1]s <path-to-template> <output-path>
e.g.,     %[1]s templates/main.gohtml templates/main.html

`, appName, version)
}
