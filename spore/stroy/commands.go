package main

var Commands = map[string][]string{
	"install": {
		`go install -v %ldflags ./spore/.`,
	},
	"run": {
		`go install -v %ldflags ./spore/.`,
		`spore`,
	},
	"hello": {
		`go install -v %ldflags ./spore/.`,
		`spore github.com/l0k18/sporeOS/cmd/hello test`,
	},
	"windows": {
		`go build -v -ldflags="-H windowsgui" %ldflags" ./spore/.`,
	},
	"update": {
		"go install -v %ldflags ./stroy/.",
	},
}
