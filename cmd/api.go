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
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/go-ego/re/log"
)

var cmdApi = &Command{
	UsageLine: "api [appname]",
	Short:     "auto-generate code for the ego app, Creates a ego API application",
	Long: `

`, Run: createApi,
}

func createApi(cmd *Command, args []string) int {
	gopath := GetGOPATHs()
	fmt.Println(gopath)
	githubsrc := gopath[0] + "/src/github.com/go-ego/re/api/"
	if runtime.GOOS == "windows" {
		githubsrc = strings.Replace(githubsrc, "/", "\\", -1)
	}
	// fmt.Println("githubsrc--------", githubsrc)

	afilesrc, err := WalkFile(githubsrc, "")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(afilesrc)

	if len(args) != 1 {
		logger.Fatal("Argument [appname] is missing")
	}

	apppath, packpath, err := checkEnv(args[0])
	if err != nil {
		logger.Fatalf("%s", err)
	}

	if isExist(apppath) {
		logger.Errorf(log.Bold("Application '%s' already exists"), apppath)
		logger.Warn(log.Bold("Do you want to overwrite it? [Yes|No] "))
		if !askForConfirmation() {
			os.Exit(2)
		}
	}

	logger.Info("Creating application... " + packpath)

	for i := 0; i < len(afilesrc); i++ {
		if runtime.GOOS == "windows" {
			afilesrc[i] = strings.Replace(afilesrc[i], "/", "\\", -1)
			apppath = strings.Replace(apppath, "/", "\\", -1)
		}
		tfile := strings.Replace(afilesrc[i], githubsrc, "", -1)
		var name string
		if runtime.GOOS == "windows" {
			name = apppath + "\\" + tfile
		} else {
			name = apppath + "/" + tfile
		}
		CopyFile(afilesrc[i], name)
	}

	return 0
}
