// Copyright 2013 Ego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/go-ego/re/log"
)

// Go is a basic promise implementation: it wraps calls a function in a goroutine
// and returns a channel which will later return the function's return value.
func Go(f func() error) chan error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()
	return ch
}

// IsExist returns whether a file or directory exists.
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// GetGOPATHs returns all paths in GOPATH variable.
func GetGOPATHs() []string {
	gopath := os.Getenv("GOPATH")
	var paths []string
	if runtime.GOOS == "windows" {
		gopath = strings.Replace(gopath, "\\", "/", -1)
		paths = strings.Split(gopath, ";")
	} else {
		paths = strings.Split(gopath, ":")
	}
	return paths
}

// IsInGOPATH checks the path is in the fisrt GOPATH(/src) or not
func IsInGOPATH(thePath string) bool {
	if runtime.GOOS == "windows" {
		thePath = filepath.ToSlash(thePath)
	}
	return strings.Contains(thePath, GetGOPATHs()[0]+"/src")
}

// IsEgoProject checks whether the current path is a Egogo application or not
func IsEgoProject(thePath string) bool {
	mainFiles := []string{}
	hasEgogoRegex := regexp.MustCompile(`(?s)package main.*?import.*?\(.*?github.com/go-ego/ego".*?\).*func main()`)
	c := make(chan error)
	// Walk the application path tree to look for main files.
	// Main files must satisfy the 'hasEgogoRegex' regular expression.
	go func() {
		filepath.Walk(thePath, func(fpath string, f os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			// Skip sub-directories
			if !f.IsDir() {
				var data []byte
				data, err = ioutil.ReadFile(fpath)
				if err != nil {
					c <- err
					return nil
				}

				if len(hasEgogoRegex.Find(data)) > 0 {
					mainFiles = append(mainFiles, fpath)
				}
			}
			return nil
		})
		close(c)
	}()

	if err := <-c; err != nil {
		logger.Fatalf("Unable to walk '%s' tree: %s", thePath, err)
	}

	if len(mainFiles) > 0 {
		return true
	}
	return false
}

// SearchGOPATHs searchs the user GOPATH(s) for the specified application name.
// It returns a boolean, the application's GOPATH and its full path.
func SearchGOPATHs(app string) (bool, string, string) {
	gps := GetGOPATHs()
	if len(gps) == 0 {
		logger.Fatal("GOPATH environment variable is not set or empty")
	}

	// Lookup the application inside the user workspace(s)
	for _, gopath := range gps {
		var currentPath string

		if !strings.Contains(app, "src") {
			gopathsrc := path.Join(gopath, "src")
			currentPath = path.Join(gopathsrc, app)
		} else {
			currentPath = app
		}

		if isExist(currentPath) {
			return true, gopath, currentPath
		}
	}
	return false, "", ""
}

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		logger.Fatalf("%s", err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

func containsString(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}

// snake string, XxYy to xx_yy
func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// camelCase converts a _ delimited string to camel case
// e.g. very_important_person => VeryImportantPerson
func camelCase(in string) string {
	tokens := strings.Split(in, "_")
	for i := range tokens {
		tokens[i] = strings.Title(strings.Trim(tokens[i], " "))
	}
	return strings.Join(tokens, "")
}

// formatSourceCode formats source files
func formatSourceCode(filename string) {
	cmd := exec.Command("gofmt", "-w", filename)
	if err := cmd.Run(); err != nil {
		logger.Warnf("Error while running gofmt: %s", err)
	}
}

// The string flag list, implemented flag.Value interface
type strFlags []string

func (s *strFlags) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *strFlags) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// CloseFile attempts to close the passed file
// or panics with the actual error
func CloseFile(f *os.File) {
	err := f.Close()
	MustCheck(err)
}

func exitPrint(con string) {
	fmt.Fprintln(os.Stderr, con)
	os.Exit(2)
}

// WriteToFile creates a file and writes content to it
func WriteToFile(filename, content string) {
	f, err := os.Create(filename)
	MustCheck(err)
	defer CloseFile(f)
	_, err = f.WriteString(content)
	MustCheck(err)
}

// Now returns the current local time in the specified layout
func Now(layout string) string {
	return time.Now().Format(layout)
}

// MustCheck panics when the error is not nil
func MustCheck(err error) {
	if err != nil {
		panic(err)
	}
}

// __FILE__ returns the file name in which the function was invoked
func __FILE__() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

// __LINE__ returns the line number at which the function was invoked
func __LINE__() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

// EgoFuncMap returns a FuncMap of functions used in different templates.
func EgoFuncMap() template.FuncMap {
	return template.FuncMap{
		"trim":       strings.TrimSpace,
		"bold":       log.Bold,
		"headline":   log.MagentaBold,
		"foldername": log.RedBold,
		"endline":    log.EndLine,
		"tmpltostr":  TmplToString,
	}
}

// TmplToString parses a text template and return the result as a string.
func TmplToString(tmpl string, data interface{}) string {
	t := template.New("tmpl").Funcs(EgoFuncMap())
	template.Must(t.Parse(tmpl))

	var doc bytes.Buffer
	err := t.Execute(&doc, data)
	MustCheck(err)

	return doc.String()
}

func checkEnv(appname string) (apppath, packpath string, err error) {
	gps := GetGOPATHs()
	if len(gps) == 0 {
		logger.Fatal("GOPATH environment variable is not set or empty")
	}
	currpath, _ := os.Getwd()
	currpath = path.Join(currpath, appname)
	for _, gpath := range gps {
		gsrcpath := path.Join(gpath, "src")
		if strings.HasPrefix(currpath, gsrcpath) {
			packpath = strings.Replace(currpath[len(gsrcpath)+1:], string(filepath.Separator), "/", -1)
			return currpath, packpath, nil
		}
	}

	// In case of multiple paths in the GOPATH, by default
	// we use the first path
	gopath := gps[0]

	logger.Warn("You current workdir is not inside $GOPATH/src.")
	logger.Debugf("GOPATH: %s", __FILE__(), __LINE__(), gopath)

	gosrcpath := path.Join(gopath, "src")
	apppath = path.Join(gosrcpath, appname)

	if _, e := os.Stat(apppath); os.IsNotExist(e) == false {
		err = fmt.Errorf("Cannot create application without removing '%s' first.", apppath)
		logger.Errorf("Path '%s' already exists", apppath)
		return
	}
	packpath = strings.Join(strings.Split(apppath[len(gosrcpath)+1:], string(filepath.Separator)), "/")
	return
}
