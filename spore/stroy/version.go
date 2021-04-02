package main

import "fmt"

var (
	URL       string
	GitRef    string
	GitCommit string
	BuildTime string
	Tag       string
)

func GetVersion() string {
	return fmt.Sprintf(
		"app information: repo: %s branch: %s commit: %s built"+
			": %s tag: %s...\n", URL, GitRef, GitCommit, BuildTime, Tag,
	)
}

var ldFlags []string
