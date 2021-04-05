package main

var commands = map[string][]string{
	"release": {
		"echo nothing",
	},
	"tests": {
		"go test ./...",
	},
	"builder": {
		"go install -v ./qubuild/.",
	},
}
