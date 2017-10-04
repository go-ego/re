// Copyright 2013 bee authors
//
// Copyright 2016 ego authors
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

// ego is a tool for developling applications based on ego framework.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/go-ego/re/cmd"
)

const (
	version string = "v0.10.0.63, Nile River!"
)

//GetVersion get version
func GetVersion() string {
	return version
}

func main() {
	// currentpath, _ := os.Getwd()

	flag.Usage = cmd.Usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()

	if len(args) < 1 {
		cmd.Usage()
	}

	if args[0] == "help" {
		cmd.Help(args[1:])
		return
	}

	for _, cmd := range cmd.AvailableCommands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}

			if cmd.PreRun != nil {
				cmd.PreRun(cmd, args)
			}

			// Check if current directory is inside the GOPATH,
			// if so parse the packages inside it.
			// if strings.Contains(currentpath, GetGOPATHs()[0]+"/src") && isGenerateDocs(cmd.Name(), args) {
			// if cmd.IsInGOPATH(currentpath) && cmd.IsGenerateDocs(cmd.Name(), args) {
			// 	parsePackagesFromDir(currentpath)
			// }

			os.Exit(cmd.Run(cmd, args))
			return
		}
	}

	cmd.PrintErrorAndExit("Unknown subcommand")
}
