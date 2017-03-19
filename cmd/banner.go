// Copyright 2017 The go-ego Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-ego/ego/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"text/template"
)

// RuntimeInfo holds information about the current runtime.
type RuntimeInfo struct {
	GoVersion  string
	GOOS       string
	GOARCH     string
	NumCPU     int
	GOPATH     string
	GOROOT     string
	Compiler   string
	EVersion   string
	EgoVersion string
}

// InitBanner loads the banner and prints it to output
// All errors are ignored, the application will not
// print the banner in case of error.
func InitBanner(out io.Writer, in io.Reader) {
	if in == nil {
		logger.Fatal("The input is nil")
	}

	banner, err := ioutil.ReadAll(in)
	if err != nil {
		logger.Fatalf("Error while trying to read the banner: %s", err)
	}

	show(out, string(banner))
}

func show(out io.Writer, content string) {
	t, err := template.New("banner").
		Funcs(template.FuncMap{"Now": Now}).
		Parse(content)

	if err != nil {
		logger.Fatalf("Cannot parse the banner template: %s", err)
	}

	err = t.Execute(out, RuntimeInfo{
		getGoVersion(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		os.Getenv("GOPATH"),
		runtime.GOROOT(),
		runtime.Compiler,
		version,
		getEgoVersion(),
		// getBeegoVersion(),
	})
	MustCheck(err)
}
