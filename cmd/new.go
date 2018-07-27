// Copyright 2016 The go-ego Project Developers. See the COPYRIGHT
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
	"log"
	"os"
	"runtime"
	"strings"

	rlog "github.com/go-ego/re/log"
	"github.com/go-vgo/gt/file"
)

var cmdNew = &Command{
	UsageLine: "new [appname]",
	Short:     "auto-generate code for the ego app, Creates a ego API application",
	Long: `

`, Run: createDir,
}

func createDir(cmd *Command, args []string) int {
	gopath := GetGOPATHs()
	log.Println("gopath: ", gopath)

	// githubsrc := gopath[0] + "/src/github.com/go-ego/re/gen/"
	githubsrc := hasFile(gopath, "/src/github.com/go-ego/re/gen/")
	newDir(githubsrc, args)

	return 0
}

func hasFile(gopath []string, name string) string {
	for i := 0; i < len(gopath); i++ {
		filename := gopath[i] + name
		if file.Exist(filename) {
			return filename
		}
	}

	return ""
}

func newDir(githubsrc string, args []string) {
	if runtime.GOOS == "windows" {
		githubsrc = strings.Replace(githubsrc, "/", "\\", -1)
	}
	// fmt.Println("githubsrc--------", githubsrc)

	filesrc, err := file.Walk(githubsrc, "")
	if err != nil {
		log.Println("walk flie: ", err)
	}
	// log.Println(filesrc)

	if len(args) != 1 {
		logger.Fatal("Argument [appname] is missing")
	}

	appPath, packPath, err := checkEnv(args[0])
	if err != nil {
		logger.Fatalf("%s", err)
	}

	if isExist(appPath) {
		logger.Errorf(rlog.Bold("Application '%s' already exists"), appPath)
		logger.Warn(rlog.Bold("Do you want to overwrite it? [Yes|No] "))
		if !askForConfirmation() {
			os.Exit(2)
		}
	}

	logger.Info("Creating application... " + packPath)

	for i := 0; i < len(filesrc); i++ {
		if runtime.GOOS == "windows" {
			filesrc[i] = strings.Replace(filesrc[i], "/", "\\", -1)
			appPath = strings.Replace(appPath, "/", "\\", -1)
		}

		tfile := strings.Replace(filesrc[i], githubsrc, "", -1)
		var name string
		if runtime.GOOS == "windows" {
			name = appPath + "\\" + tfile
		} else {
			name = appPath + "/" + tfile
		}

		CopyFile(filesrc[i], name)
	}
}

func CopyFile(src, dst string) {
	if !file.Exist(dst) {
		Writefile(dst, "")
	}

	var redst string
	if runtime.GOOS == "windows" {
		if strings.Contains(dst, ".") {
			dstarr := strings.Split(dst, "\\")
			len := len(dstarr) - 1
			datstr := dstarr[len]
			redst = strings.Replace(dst, datstr, "", -1)
		} else {
			redst = dst
		}
		os.MkdirAll(redst, os.ModePerm)
	}

	file.CopyFile(src, dst)
}

func Writefile(fileName, writeStr string) {
	log.Println(rlog.Blue("Create:::"), rlog.Yellow(fileName))
	file.Write(fileName, writeStr)
}
