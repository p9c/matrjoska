package main

import (
	pod2 "github.com/p9c/monorepo/pkg/pod"
	"github.com/p9c/monorepo/pkg/podopts"
	"github.com/p9c/monorepo/pod"
	"github.com/p9c/monorepo/version"
)

func main() {
	I.Ln(version.Get())
	var cx *pod2.State
	var e error
	if cx, e = pod.GetNewContext(podopts.GetDefaultConfig()); E.Chk(e) {
	}
	_ = cx
	I.S(cx)
}
