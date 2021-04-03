package version

import "fmt"

var (
	
	// URL is the git URL for the repository
	URL = "github.com/p9c/monorepo"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "refs/heads/main"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "91786c6dd9775a933968b3aa8fba977a28a80990"
	// BuildTime stores the time when the current binary was built
	BuildTime = "2021-03-28T03:41:47+02:00"
	// Tag lists the Tag on the build, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "v0.0.2"
	// PathBase is the path base returned from runtime caller
	PathBase = "/home/loki/src/github.com/p9c/monorepo/spore/"
)

// Get returns a pretty printed version information string
func Get() string {
	return fmt.Sprint(
		"Build information\n"+
			"	git repository: "+URL+"\n",
		"	branch: "+GitRef+"\n"+
			"	commit: "+GitCommit+"\n"+
			"	built: "+BuildTime+"\n"+
			"	Tag: "+Tag+"\n",
	)
}
