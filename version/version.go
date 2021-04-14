package version

import "fmt"

var (

	// URL is the git URL for the repository
	URL = "github.com/p9c/monorepo"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "refs/heads/main"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "8d89a3d3d2881f77a3bc72db8f481966f6df6b9a"
	// BuildTime stores the time when the current binary was built
	BuildTime = "2021-04-14T09:52:05+02:00"
	// Tag lists the Tag on the podbuild, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "v0.0.13+"
	// PathBase is the path base returned from runtime caller
	PathBase = "/home/loki/src/github.com/p9c/monorepo/"
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
