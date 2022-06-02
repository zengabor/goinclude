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
	version = "1.1"
)

var re = regexp.MustCompile(`{{\s*include\s+"([^"]+)"\s*}}`)

func main() {
	switch {
	case len(os.Args) == 1:
		help()
	case len(os.Args) == 3:
		mustIncludeFiles(os.Args[1], os.Args[2])
	default:
		log.Fatalf("%s: provide input file path and output file path as parameters", appName)
	}
}

func mustIncludeFiles(pathToTemplate, pathToOutputFile string) {
	tmpl, err := getFileContent(pathToTemplate)
	if err != nil {
		log.Fatalf("%s: could not read %q", appName, pathToTemplate)
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
	var errors []string
	out := re.ReplaceAllStringFunc(tmpl, func(s string) string {
		toBeIncludedFilePath := re.ReplaceAllString(s, "$1")
		c, err := get(toBeIncludedFilePath)
		if err != nil {
			e := fmt.Sprintf("error in %s: %s", pathToTemplate, err)
			errors = append(errors, e)
			return fmt.Sprintf("(%s: %s)", appName, e)
		}
		return c
	})
	if len(includedFiles) > 0 {
		fmt.Println(strings.Join(includedFiles, " "))
	}
	if err := ioutil.WriteFile(pathToOutputFile, []byte(out), 0644); err != nil {
		log.Fatalf("%s: error writing %q", appName, pathToOutputFile)
	}
	if len(errors) > 0 {
		log.Fatalf("%s: \n - %v", appName, strings.Join(errors, "\n - "))
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
			log.Fatalf("%s: getting working directory: %s\n", appName, err)
		}
		pathToFile = path.Join(workingDirectory, pathToFile)
	}
	s, err := filepath.Abs(pathToFile)
	if err != nil {
		log.Fatalf("%s: resolving full path: %s\n", appName, err)
	}
	return s
}

func help() {
	fmt.Printf(`%s %s // github.com/zengabor/goinclude
Includes the content of the indicated files into a go template and outputs the merged result to the specified file.

This is how you include a file inside a template: {{ include "path/otherfile.html" }}

  Usage:  %[1]s <path-to-template> <output-path>
  E.g.,   %[1]s templates/main.gohtml templates/main.html

`, appName, version)
}
