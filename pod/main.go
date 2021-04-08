package main

import (
	"bytes"
	"encoding/json"
	"github.com/p9c/monorepo/pkg/log"
	"github.com/p9c/monorepo/pkg/pod"
	"github.com/p9c/monorepo/pkg/podopts"
	"github.com/p9c/monorepo/pod/launchers"
	"github.com/p9c/monorepo/pod/podcfgs"
	"github.com/p9c/monorepo/version"
	"os"
)

func main() {
	log.SetLogLevel("trace")
	I.Ln(version.Get())
	var cx *pod.State
	var e error
	if cx, e = launchers.GetNewContext(podcfgs.GetDefaultConfig()); E.Chk(e) {
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
