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
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-ego/e/log"
)

var cmdNew = &Command{
	UsageLine: "new [appname]",
	Short:     "auto-generate code for the ego app, Creates a ego API application",
	Long: `

`, Run: createDir,
}

func createDir(cmd *Command, args []string) int {
	gopath := GetGOPATHs()
	fmt.Println(gopath)
	githubsrc := gopath[0] + "/src/github.com/go-ego/e/gen/"
	if runtime.GOOS == "windows" {
		githubsrc = strings.Replace(githubsrc, "/", "\\", -1)
	}
	// fmt.Println("githubsrc--------", githubsrc)

	afile, err := WalkFile(githubsrc, "")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(afile)

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

	for i := 0; i < len(afile); i++ {

		tfile := strings.Replace(afile[i], githubsrc, "", -1)
		name := apppath + "/" + tfile

		CopyFile(afile[i], name)
	}

	return 0
}

func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	// if fileExist(dst) != true {
	if !fileExist(dst) {
		Wirtefile("", dst)
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

func Wirtefile(wirtestr string, userFile string) {

	fmt.Println(log.Blue("Create:::"), log.Yellow(userFile))

	os.MkdirAll(path.Dir(userFile), os.ModePerm)

	fout, err := os.Create(userFile)
	defer fout.Close()
	if err != nil {
		fmt.Println(userFile, err)
		return
	}

	fout.WriteString(wirtestr)

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
