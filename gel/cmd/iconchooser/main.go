package main

import (
	l "gioui.org/layout"
	"github.com/p9c/monorepo/gel"
	"github.com/p9c/monorepo/interrupt"
	"github.com/p9c/monorepo/qu"
)

type State struct {
	*gel.Window
}

func NewState(quit qu.C) *State {
	return &State{
		Window: gel.NewWindowP9(quit),
	}
}

func main() {
	quit := qu.T()
	state := NewState(quit)
	var e error
	rootWidget := state.rootWidget()
	if e = state.Window.
		Size(48, 32).
		Title("icons chooser").
		Open().
		Run(rootWidget,
			nil, func() {
				interrupt.Request()
				quit.Q()
			}, quit,
		); E.Chk(e) {
		
	}
}

func (s *State) rootWidget() l.Widget {
	return s.H2("icons").Fn
}
