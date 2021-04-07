package version

import "fmt"

var (

	// URL is the git URL for the repository
	URL = "github.com/p9c/monorepo"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "refs/heads/main"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "35c5452b6dccd9b8d67409ac2ca9ef8b5529ae52"
	// BuildTime stores the time when the current binary was built
	BuildTime = "2021-04-07T20:26:21+02:00"
	// Tag lists the Tag on the podbuild, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "v0.0.10+"
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
