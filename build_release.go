//go:build ignore_build_go
// +build ignore_build_go

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/nimaaskarian/ydo/utils"
)

func main() {
	dry := false
	ver := os.Args[1]
	if len(os.Args) >= 3 && os.Args[2] == "list" {
		dry = true
	}
	exec.Command("go", "test", "./...")
	tags := [][]string{
		[]string{
			"jalali", "webgui",
		},
		[]string{
			"jalali",
		},
		[]string{
			"webgui",
		},
		[]string{},
	}
	for _, env := range envs {
		if err := build(env, tags, ver, dry); err != nil {
			log.Fatal(err)
		}
	}
}

type Env map[string]string

var envs = []Env{
	Env{
		"GOARCH": "amd64",
		"GOOS":   "linux",
	},
	Env{
		"GOARCH": "amd64",
		"GOOS":   "windows",
	},
	Env{
		"GOARCH":      "arm64",
		"CC":          "/opt/android-sdk/ndk/27.0.12077973//toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang",
		"CGO_ENABLED": "1",
		"GOOS":        "android",
	},
}

func build(env Env, tags_array [][]string, ver string, dry bool) error {
	for _, tags := range tags_array {
		tag := strings.Join(tags, ",")
		tag_print := ""
		if len(tags) > 0 {
			tag_print = fmt.Sprintf("-%s", strings.Join(tags, "-"))
		}
		format := ""
		if env["GOOS"] == "windows" {
			format = ".exe"
		}
		name := fmt.Sprintf("ydo_%s_%s_%s%s%s", ver, env["GOOS"], env["GOARCH"], tag_print, format)
		if dry {
			fmt.Println(name)
			continue
		}
		var cmd *exec.Cmd
		if tag != "" {
			cmd = exec.Command("go", "build", "-o", name, "-tags", tag)
		} else {
			cmd = exec.Command("go", "build", "-o", name)
		}
		cmd.Env = os.Environ()
		for k, v := range env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
		utils.CmdStdOs(cmd)
		fmt.Println(cmd.String())
		err := cmd.Run()
		cmd.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}
