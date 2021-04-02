package spore

import (
	"os"
)

func (s *Shell) Main() int {
	if len(os.Args) > 1 {
		c := NewCLI(s)
		return c.Run()
	} else {
		g := NewGUI(s)
		return g.Run()
	}
}
