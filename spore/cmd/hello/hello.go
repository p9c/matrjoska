package main

import (
	"os"
	
	. "github.com/l0k18/sporeOS/pkg/log"
)

func main() {
	Debug("hello", os.Args[1:])
}
