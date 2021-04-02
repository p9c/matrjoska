package version

import "fmt"

var (

	// URL is the git URL for the repository
	URL = "github.com/p9c/duod"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "refs/heads/master"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "b512313e6a14903710f0cb3567215859960496f9"
	// BuildTime stores the time when the current binary was built
	BuildTime = "2021-04-02T20:14:57+02:00"
	// Tag lists the Tag on the build, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "+"
	// PathBase is the path base returned from runtime caller
	PathBase = "/home/loki/src/github.com/p9c/glom/vendor/github.com/p9c/gel/gelbuild/stroy.go/duod/"
)

// Get returns a pretty printed version information string
func Get() string {
	return fmt.Sprint(
		"ParallelCoin Pod\n"+
		"	git repository: "+URL+"\n",
		"	branch: "+GitRef+"\n"+
		"	commit: "+GitCommit+"\n"+
		"	built: "+BuildTime+"\n"+
		"	Tag: "+Tag+"\n",
	)
}
