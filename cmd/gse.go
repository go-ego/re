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
	"log"
)

var cmdGse = &Command{
	UsageLine: "gse [appname]",
	Short:     "auto-generate code for the gse application",
	Long: `

`, Run: createGse,
}

func createGse(cmd *Command, args []string) int {
	gopath := GetGOPATHs()
	log.Println("gopath: ", gopath)

	githubsrc := gopath[0] + "/src/github.com/go-ego/gse/data"
	newDir(githubsrc, args)

	return 0
}
