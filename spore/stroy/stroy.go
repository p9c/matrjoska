// +build !windows

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	
	"github.com/l0k18/sporeOS/pkg/log"
	"github.com/l0k18/sporeOS/pkg/util"
)

func main() {
	fmt.Println(GetVersion())
	var err error
	var ok bool
	var home string
	if home, ok = os.LookupEnv("HOME"); !ok {
		panic(err)
	}
	if len(os.Args) > 1 {
		folderName := "test0"
		var dataDir string
		if len(os.Args) > 2 {
			dataDir = os.Args[2]
		} else {
			dataDir = filepath.Join(home, folderName)
		}
		if list, ok := Commands[os.Args[1]]; ok {
			populateVersionFlags()
			for i := range list {
				// inject the data directory
				var split []string
				out := strings.ReplaceAll(list[i], "%datadir", dataDir)
				split = strings.Split(out, " ")
				for i := range split {
					split[i] = strings.ReplaceAll(
						split[i], "%ldflags",
						fmt.Sprintf(
							`-ldflags=%s`, strings.Join(
								ldFlags,
								" ",
							),
						),
					)
				}
				fmt.Printf(
					`executing item %d of list '%v' '%v' '%v'

`, i, os.Args[1],
					split[0], split[1:],
				)
				var cmd *exec.Cmd
				scriptPath := filepath.Join(util.Dir("stroy", false), "stroy.sh")
				util.EnsureDir(scriptPath)
				if err = ioutil.WriteFile(
					scriptPath,
					[]byte(strings.Join(split, " ")),
					0700,
				); log.Check(err) {
				} else {
					cmd = exec.Command("sh", scriptPath)
					cmd.Stdout = os.Stdout
					cmd.Stdin = os.Stdin
					cmd.Stderr = os.Stderr
				}
				if cmd == nil {
					panic("cmd is nil")
				}
				if err := cmd.Start(); log.Check(err) {
					os.Exit(1)
				}
				if err := cmd.Wait(); log.Check(err) {
					os.Exit(1)
				}
			}
		} else {
			fmt.Println("command", os.Args[1], "not found")
		}
	} else {
		fmt.Println("no command requested, available:")
		for i := range Commands {
			fmt.Println(i)
			for j := range Commands[i] {
				fmt.Println("\t" + Commands[i][j])
			}
		}
	}
}
