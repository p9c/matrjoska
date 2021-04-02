package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	
	. "github.com/l0k18/sporeOS/pkg/log"
	
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

func populateVersionFlags() bool {
	BuildTime = time.Now().Format(time.RFC3339)
	var cwd string
	var err error
	if cwd, err = os.Getwd(); Check(err) {
		return false
	}
	var repo *git.Repository
	if repo, err = git.PlainOpen(cwd); Check(err) {
		return false
	}
	var rr []*git.Remote
	if rr, err = repo.Remotes(); Check(err) {
		return false
	}
	// spew.Dump(rr)
	for i := range rr {
		rs := rr[i].String()
		if strings.HasPrefix(rs, "origin") {
			rss := strings.Split(rs, "git@")
			if len(rss) > 1 {
				rsss := strings.Split(rss[1], ".git")
				URL = strings.ReplaceAll(rsss[0], ":", "/")
				break
			}
			rss = strings.Split(rs, "https://")
			if len(rss) > 1 {
				rsss := strings.Split(rss[1], ".git")
				URL = rsss[0]
				break
			}
			
		}
	}
	var rh *plumbing.Reference
	if rh, err = repo.Head(); Check(err) {
		return false
	}
	rhs := rh.Strings()
	GitRef = rhs[0]
	GitCommit = rhs[1]
	var rt storer.ReferenceIter
	if rt, err = repo.Tags(); Check(err) {
		return false
	}
	var maxVersion int
	var maxString string
	var maxIs bool
	if err = rt.ForEach(
		func(pr *plumbing.Reference) error {
			prs := strings.Split(pr.String(), "/")[2]
			if strings.HasPrefix(prs, "v") {
				var va [3]int
				_, _ = fmt.Sscanf(prs, "v%d.%d.%d", &va[0], &va[1], &va[2])
				vn := va[0]*1000000 + va[1]*1000 + va[2]
				if maxVersion < vn {
					maxVersion = vn
					maxString = prs
				}
				if pr.Hash() == rh.Hash() {
					maxIs = true
				}
			}
			return nil
		},
	); Check(err) {
		return false
	}
	if !maxIs {
		maxString += "+"
	}
	Tag = maxString
	ldFlags = []string{
		`"-X main.URL=` + URL + ``,
		`-X main.GitCommit=` + GitCommit + ``,
		`-X main.BuildTime=` + BuildTime + ``,
		`-X main.GitRef=` + GitRef + ``,
		`-X main.Tag=` + Tag + `"`,
	}
	return true
}
