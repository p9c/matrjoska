package main

var commands = map[string][]string{
	"build": {
		"go build -v ./...",
	},
	"install": {
		"go install -v ./...",
	},
	"tidy": {
		"go mod tidy",
	},
}
