package main

var commands = map[string][]string{
	"build": {
		"go build -v ./pod/pod/.",
	},
	"install": {
		"go install -v ./pod/pod/.",
	},
	"headless": {
		"go install -v -tags headless ./pod/pod/.",
	},
	"docker": {
		"go install -v -tags headless ./pod/pod/.",
	},
	"tests": {
		"go test ./...",
	},
	"builder": {
		"go install -v ./pod/podbuild/.",
	},
}
