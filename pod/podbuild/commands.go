package main

var commands = map[string][]string{
	"build": {
		"go build -v ./pod/.",
	},
	"install": {
		"go install -v ./pod/.",
	},
	"run": {
		"go install -v ./pod/.",
		"pod",
	},
	"node": {
		"go install -v ./pod/.",
		"pod node",
	},
	"headless": {
		"go install -v -tags headless ./pod/.",
	},
	"docker": {
		"go install -v -tags headless ./pod/.",
	},
	"tests": {
		"go test ./...",
	},
	"builder": {
		"go install -v ./pod/podbuild/.",
	},
}
