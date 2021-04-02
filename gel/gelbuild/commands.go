package main

var commands = map[string][]string{
	"build": {
		"go build -v",
	},
	"tests": {
		"go test ./...",
	},
	"builder": {
		"go install -v ./gelbuild/.",
	},
}
