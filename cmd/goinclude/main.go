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
	case len(os.Args) < 2:
		help()
	case len(os.Args) == 2:
		includeFiles(os.Args[1])
	default:
		log.Fatal(appName + ": provide exactly one file path as parameter")
	}
}

func includeFiles(pathToTemplate string) {
	tmpl, err := getFileContent(pathToTemplate)
	if err != nil {
		log.Fatal(appName + ": could not open " + pathToTemplate)
	}
	fmt.Fprint(os.Stdout, re.ReplaceAllStringFunc(tmpl, func(s string) string {
		c, err := getFileContent(re.ReplaceAllString(s, "$1"))
		if err != nil {
			return fmt.Sprintf(
				"(%s: error including %s into %s: %s)",
				appName, getFullPath(pathToTemplate), pathToTemplate, err,
			)
		}
		return c
	}))
}

func getFileContent(filePath string) (s string, err error) {
	b, err := ioutil.ReadFile(getFullPath(filePath))
	if err != nil {
		return
	}
	return string(b), nil
}

func getFullPath(pathToFile string) string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(fmt.Sprintf("%s: could not get working directory: %s\n", appName, err))
	}
	if !strings.HasPrefix(pathToFile, "/") {
		pathToFile = path.Join(workingDirectory, pathToFile)
	}
	s, err := filepath.Abs(pathToFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s: could not resolve full path: %s\n", appName, err))
	}
	return s
}

func help() {
	fmt.Printf(`%s %s // github.com/zengabor/goinclude
Includes the content of a file into a go template and outputs the merged result to stdout.

Usage:    %[1]s <pathToTemplateFile>

Example:
  %[1]s templates/main.gohtml > templates/main.html

`, appName, version)
}
