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
)

func Colorize(text string, status string) string {
	out := ""
	switch status {
	case "succ":
		out = "\033[32;1m" // Blue
	case "fail":
		out = "\033[31;1m" // Red
	case "warn":
		out = "\033[33;1m" // Yellow
	case "note":
		out = "\033[34;1m" // Green
	case "blue":
		out = "\033[44;1m" // blue
	default:
		out = "\033[0m" // Default
	}
	return out + text + "\033[0m"
}

var cmdFix = &Command{
	UsageLine: "fix",
	Short:     "fix the ego application to make it compatible with ego 1.0",
	Long: `
As from ego 1.0, there's some incompatible code with the old version.

bee fix help to upgrade the application to ego 1.0
`,
}

func init() {
	cmdFix.Run = runFix
}

func runFix(cmd *Command, args []string) int {
	fmt.Println(Colorize("There is no fix in this version!", "note"))
	return 0
}
