package version

import "fmt"

var (

	// URL is the git URL for the repository
	URL = "github.com/p9c/log"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "refs/heads/main"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "78568bfd404f75f1e4f5796e3b99906a255269ff"
	// BuildTime stores the time when the current binary was built
	BuildTime = "2021-04-02T22:21:18+02:00"
	// Tag lists the Tag on the build, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "v0.0.6"
	// PathBase is the path base returned from runtime caller
	PathBase = "/home/loki/src/github.com/p9c/log/"
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
