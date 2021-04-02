package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	
	. "github.com/l0k18/sporeOS/pkg/log"
)

var WindowsExec = func(split []string) (out *exec.Cmd) {
	out = exec.Command(split[0])
	out.SysProcAttr = &syscall.SysProcAttr{}
	out.SysProcAttr.CmdLine = strings.Join(split, " ")
	return
}

func main() {
	fmt.Println(GetVersion())
	var err error
	var ok bool
	var home string
	if runtime.GOOS == "windows" {
		var homedrive string
		if homedrive, ok = os.LookupEnv("HOMEDRIVE"); !ok {
			panic(err)
		}
		var homepath string
		if homepath, ok = os.LookupEnv("HOMEPATH"); !ok {
			panic(err)
		}
		home = homedrive + homepath
	} else {
		if home, ok = os.LookupEnv("HOME"); !ok {
			panic(err)
		}
	}
	if len(os.Args) > 1 {
		folderName := "test0"
		var datadir string
		if len(os.Args) > 2 {
			datadir = os.Args[2]
		} else {
			datadir = filepath.Join(home, folderName)
		}
		if list, ok := Commands[os.Args[1]]; ok {
			populateVersionFlags()
			for i := range list {
				var split []string
				out := strings.ReplaceAll(list[i], "%datadir", datadir)
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
				cmd := WindowsExec(split)
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				if err := cmd.Start(); Check(err) {
					Infos(err)
					os.Exit(1)
				}
				if err := cmd.Wait(); Check(err) {
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
