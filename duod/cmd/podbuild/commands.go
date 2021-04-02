package main

var commands = map[string][]string{
	"build": {
		"go build -v ./...",
	},
	"test": {
		"go test ./...",
	},
	"gen": {
		"go generate ./...",
	},
	"install": {
		"go install -v",
	},
	"headless": {
		"go install -v -tags headless",
	},
	"builder": {
		"go install -v ./cmd/podbuild/.",
	},
}
