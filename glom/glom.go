package main

import (
	"github.com/p9c/monorepo/gel"
	"github.com/p9c/monorepo/glom/pkg/pathtree"
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
	folderView := pathtree.New(state.Window)
	if e = state.Window.
		Size(20, 20).
		Title("glom, the visual code editor").
		Open().
		Run(folderView.Fn,
			nil, func() {
				interrupt.Request()
				quit.Q()
			}, quit,
		); E.Chk(e) {
		
	}
}
