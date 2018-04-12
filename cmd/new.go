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
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	rlog "github.com/go-ego/re/log"
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

	githubsrc := gopath[0] + "/src/github.com/go-ego/re/gen/"
	newDir(githubsrc, args)

	return 0
}

func newDir(githubsrc string, args []string) {
	if runtime.GOOS == "windows" {
		githubsrc = strings.Replace(githubsrc, "/", "\\", -1)
	}
	// fmt.Println("githubsrc--------", githubsrc)

	filesrc, err := WalkFile(githubsrc, "")
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

func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

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
	// if fileExist(dst) != true {
	if !fileExist(dst) {
		Writefile("", dst)
	}
	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func Writefile(writeStr string, userFile string) {

	log.Println(rlog.Blue("Create:::"), rlog.Yellow(userFile))

	os.MkdirAll(path.Dir(userFile), os.ModePerm)

	fout, err := os.Create(userFile)
	defer fout.Close()
	if err != nil {
		fmt.Println(userFile, err)
		return
	}

	fout.WriteString(writeStr)

}

func WalkFile(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {

		if fi.IsDir() { // dir
			return nil
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})

	return files, err
}
