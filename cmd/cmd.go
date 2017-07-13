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
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/go-ego/re/log"
)

const (
	version string = "v0.10.0.41, Nile River!"
)

// Command is the unit of execution
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string) int

	// PreRun performs an operation before running the command
	PreRun func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short string

	// Long is the long message shown in the 'go help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool

	// output out writer if set in SetOutput(w)
	output *io.Writer
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (c *Command) SetOutput(output io.Writer) {
	c.output = &output
}

// Out returns the out writer of the current command.
// If cmd.output is nil, os.Stderr is used.
func (c *Command) Out() io.Writer {
	if c.output != nil {
		return *c.output
	}
	return log.NewColorWriter(os.Stderr)
}

// Usage puts out the usage for the command.
func (c *Command) Usage() {
	tmpl(cmdUsage, c)
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as import path.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (c *Command) Options() map[string]string {
	options := make(map[string]string)
	c.Flag.VisitAll(func(f *flag.Flag) {
		defaultVal := f.DefValue
		if len(defaultVal) > 0 {
			// if strings.Contains(defaultVal, ":") {
			// 	// Truncate the flag's default value by appending '...' at the end
			// 	options[f.Name+"="+strings.Split(defaultVal, ":")[0]+":..."] = f.Usage
			// } else {
			options[f.Name+"="+defaultVal] = f.Usage
			// }
		} else {
			options[f.Name] = f.Usage
		}
	})
	return options
}

var AvailableCommands = []*Command{
	// cmdGen,
	cmdNew,
	cmdRun,
	cmdPack,
	cmdApi,
	cmdRiot,
	//cmdRouter,
	//cmdTest,
	cmdBale,
	cmdVersion,
	// cmdGenerate,
	//cmdRundocs,
	cmdMigrate,
	cmdFix,
	// cmdDockerize,
}

var logger = log.GetEgoLogger(os.Stdout)

func IsGenerateDocs(name string, args []string) bool {
	if name != "generate" {
		return false
	}
	for _, a := range args {
		if a == "docs" {
			return true
		}
	}
	return false
}

var usageTemplate = `re is a Fast and Flexible tool for managing your ego Web Application.

{{"USAGE" | headline}}
    {{"re command [arguments]" | bold}}

{{"AVAILABLE COMMANDS" | headline}}
{{range .}}{{if .Runnable}}
    {{.Name | printf "%-11s" | bold}} {{.Short}}{{end}}{{end}}

Use {{"re help [command]" | bold}} for more information about a command.

{{"ADDITIONAL HELP TOPICS" | headline}}
{{range .}}{{if not .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use {{"re help [topic]" | bold}} for more information about that topic.
`

var helpTemplate = `{{"USAGE" | headline}}
  {{.UsageLine | printf "re %s" | bold}}
{{if .Options}}{{endline}}{{"OPTIONS" | headline}}{{range $k,$v := .Options}}
  {{$k | printf "-%s" | bold}}
      {{$v}}
  {{end}}{{end}}
{{"DESCRIPTION" | headline}}
  {{tmpltostr .Long . | trim}}
`

var errorTemplate = `re: %s.
Use {{"re help" | bold}} for more information.
`

var cmdUsage = `Use {{printf "re help %s" .Name | bold}} for more information.{{endline}}`

func Usage() {
	tmpl(usageTemplate, AvailableCommands)
	os.Exit(2)
}

func tmpl(text string, data interface{}) {
	output := log.NewColorWriter(os.Stderr)

	t := template.New("usage").Funcs(EgoFuncMap())
	template.Must(t.Parse(text))

	err := t.Execute(output, data)
	MustCheck(err)
}

func Help(args []string) {
	if len(args) == 0 {
		Usage()
	}
	if len(args) != 1 {
		PrintErrorAndExit("Too many arguments")
	}

	arg := args[0]

	for _, cmd := range AvailableCommands {
		if cmd.Name() == arg {
			tmpl(helpTemplate, cmd)
			return
		}
	}
	PrintErrorAndExit("Unknown help topic")
}

func PrintErrorAndExit(message string) {
	tmpl(fmt.Sprintf(errorTemplate, message), nil)
	os.Exit(2)
}
