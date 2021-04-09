package main

import (
	"bytes"
	"encoding/json"
	_ "gioui.org/app/permission/networkstate" // todo: integrate this into routeable package
	_ "gioui.org/app/permission/storage"      // this enables the home folder appdata directory to work on android (and ios)
	"github.com/p9c/monorepo/pod/context"
	
	// _ "gioui.org/app/permission/bluetooth"
	// _ "gioui.org/app/permission/camera"
	"github.com/p9c/monorepo/pkg/log"
	"github.com/p9c/monorepo/pkg/pod"
	"github.com/p9c/monorepo/pkg/podopts"
	"github.com/p9c/monorepo/pod/podcfgs"
	"github.com/p9c/monorepo/version"
	"os"
)

func main() {
	log.SetLogLevel("info")
	I.Ln(version.Get())
	var cx *pod.State
	var e error
	if cx, e = context.GetNew(podcfgs.GetDefaultConfig()); E.Chk(e) {
		fail()
	}
	
	// if e = debugConfig(cx.Config); E.Chk(e) {
	// 	fail()
	// }
	
	D.Ln("running command", cx.Config.RunningCommand.Name)
	if e = cx.Config.RunningCommand.Entrypoint(cx); E.Chk(e) {
		fail()
	}
}

func debugConfig(c *podopts.Config) (e error) {
	c.ShowAll = true
	defer func() { c.ShowAll = false }()
	var j []byte
	if j, e = c.MarshalJSON(); E.Chk(e) {
		return
	}
	var b []byte
	jj := bytes.NewBuffer(b)
	if e = json.Indent(jj, j, "", "\t"); E.Chk(e) {
		return
	}
	I.Ln(jj.String())
	return
}

func fail() {
	os.Exit(1)
}
