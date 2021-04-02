package spore

import (
	. "github.com/l0k18/sporeOS/pkg/log"
)

type GUI struct {
	*Shell
}

func NewGUI(s *Shell) *GUI {
	return &GUI{Shell: s}
}

func (c *GUI) Run() int {
	Debug("running spore GUI")
	return 0
}
