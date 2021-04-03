package main

import (
	l "gioui.org/layout"
	"github.com/p9c/monorepo/gel"
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
		Title("hello world").
		Open().
		Run(
			rootWidget,
			nil,
			quit.Q,
			quit,
		); E.Chk(e) {
	}
}

func (s *State) rootWidget() l.Widget {
	return s.Direction().Center().Embed(s.H2("hello world!").Fn).Fn
}
