package spore

import (
	"os"
	
	. "github.com/l0k18/sporeOS/pkg/log"
)

type CLI struct {
	*Shell
}

func NewCLI(s *Shell) *CLI {
	return &CLI{Shell: s}
}

func (c *CLI) Run() int {
	Debug("running", os.Args[1:])
	Debug("dataDir", c.dataDir)
	// first we need to check we have go in the shell environment, fetch and unpack it otherwise, write config setting
	// which url and hash of package, write command to set new URL and hash for Go distribution.
	// maybe also command to list and select other than the main/master branch
	// packages will need commit/tag, and go version specifications, and multiple go version builds, to pull and cache
	// go distribution when required and not cached
	// check if we already have package and the version is current
	// storage for packages - path and then git commit hash
	// prune binary cache. keep whole git repo tied to package, and refer to cache if base is requested again
	// build cache management command - set version to build and use by git commit or tag
	// `go build github.com/org/repo`
	// move binary to binary cache
	// run binary with given parameters
	return 0
}
