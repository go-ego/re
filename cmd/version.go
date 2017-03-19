package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	path "path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/go-ego/e/log"

	"gopkg.in/yaml.v2"
)

var cmdVersion = &Command{
	UsageLine: "version",
	Short:     "Prints the current Bee version",
	Long: `
Prints the current Bee, Beego and Go version alongside the platform information.

`,
	Run: versionCmd,
}

const verboseVersionBanner string = `%s%s______
███████╗ ██████╗  ██████╗
██╔════╝██╔════╝ ██╔═══██╗
█████╗  ██║  ███╗██║   ██║
██╔══╝  ██║   ██║██║   ██║
███████╗╚██████╔╝╚██████╔╝
╚══════╝ ╚═════╝  ╚═════╝  v{{ .EVersion }}%s
%s%s
├── Ego     : {{ .EgoVersion }}
├── GoVersion : {{ .GoVersion }}
├── GOOS      : {{ .GOOS }}
├── GOARCH    : {{ .GOARCH }}
├── NumCPU    : {{ .NumCPU }}
├── GOPATH    : {{ .GOPATH }}
├── GOROOT    : {{ .GOROOT }}
├── Compiler  : {{ .Compiler }}
└── Date      : {{ Now "Monday, 2 Jan 2006" }}%s
`

const shortVersionBanner = `_________
███████╗ ██████╗  ██████╗
██╔════╝██╔════╝ ██╔═══██╗
█████╗  ██║  ███╗██║   ██║
██╔══╝  ██║   ██║██║   ██║
███████╗╚██████╔╝╚██████╔╝
╚══════╝ ╚═════╝  ╚═════╝ | v{{ .EVersion }}

---------------------------------
├── Ego     : {{ .EgoVersion }}
├── GoVersion : {{ .GoVersion }}
---------------------------------

`

var outputFormat string

func init() {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.StringVar(&outputFormat, "o", "", "Set the output format. Either json or yaml.")
	cmdVersion.Flag = *fs
}

func versionCmd(cmd *Command, args []string) int {
	cmd.Flag.Parse(args)
	stdout := cmd.Out()

	if outputFormat != "" {
		runtimeInfo := RuntimeInfo{
			getGoVersion(),
			runtime.GOOS,
			runtime.GOARCH,
			runtime.NumCPU(),
			os.Getenv("GOPATH"),
			runtime.GOROOT(),
			runtime.Compiler,
			version,
			getBeegoVersion(),
		}
		switch outputFormat {
		case "json":
			{
				b, err := json.MarshalIndent(runtimeInfo, "", "    ")
				MustCheck(err)
				fmt.Println(string(b))
				return 0
			}
		case "yaml":
			{
				b, err := yaml.Marshal(&runtimeInfo)
				MustCheck(err)
				fmt.Println(string(b))
				return 0
			}
		}
	}

	coloredBanner := fmt.Sprintf(verboseVersionBanner, "\x1b[35m", "\x1b[1m",
		"\x1b[0m", "\x1b[32m", "\x1b[1m", "\x1b[0m")
	InitBanner(stdout, bytes.NewBufferString(coloredBanner))
	return 0
}

// ShowShortVersionBanner prints the short version banner.
func ShowShortVersionBanner() {
	output := log.NewColorWriter(os.Stdout)
	InitBanner(output, bytes.NewBufferString(log.MagentaBold(shortVersionBanner)))
}

// ShowVerboseVersionBanner prints the verbose version banner
func ShowVerboseVersionBanner() {
	w := log.NewColorWriter(os.Stdout)
	coloredBanner := fmt.Sprintf(verboseVersionBanner, "\x1b[35m", "\x1b[1m", "\x1b[0m",
		"\x1b[32m", "\x1b[1m", "\x1b[0m")
	InitBanner(w, bytes.NewBufferString(coloredBanner))
}

func getEgoVersion() string {
	gopath := os.Getenv("GOPATH")
	re, err := regexp.Compile(`VERSION = "([0-9.]+)"`)
	if err != nil {
		return ""
	}
	if gopath == "" {
		err = fmt.Errorf("You need to set GOPATH environment variable")
		return ""
	}
	wgopath := path.SplitList(gopath)
	for _, wg := range wgopath {
		wg, _ = path.EvalSymlinks(path.Join(wg, "src", "github.com", "go-ego", "ego"))
		filename := path.Join(wg, "ego.go")
		_, err := os.Stat(filename)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			logger.Error("Error while getting stats of 'ego.go'")
		}
		fd, err := os.Open(filename)
		if err != nil {
			logger.Error("Error while reading 'ego.go'")
			continue
		}
		reader := bufio.NewReader(fd)
		for {
			byteLine, _, er := reader.ReadLine()
			if er != nil && er != io.EOF {
				return ""
			}
			if er == io.EOF {
				break
			}
			line := string(byteLine)
			s := re.FindStringSubmatch(line)
			if len(s) >= 2 {
				return s[1]
			}
		}

	}
	return "Ego is not installed. Please do consider installing it first: https://github.com/go-ego/ego"

}

func getGoVersion() string {
	var (
		cmdOut []byte
		err    error
	)

	if cmdOut, err = exec.Command("go", "version").Output(); err != nil {
		logger.Fatalf("There was an error running 'go version' command: %s", err)
	}
	return strings.Split(string(cmdOut), " ")[2]
}
