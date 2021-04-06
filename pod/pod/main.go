package main

import (
	"github.com/p9c/monorepo/version"
)

func main() {
	I.Ln(version.Get())
}
