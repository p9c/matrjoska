package main

import (
	l "gioui.org/layout"
	"github.com/p9c/gel"
	"github.com/p9c/interrupt"
	"github.com/p9c/qu"
)

type State struct {
	*gel.Window
}

func NewState() *State {
	return &State{}
}

func main() {
	state := NewState()
	quit := qu.T()
	var e error
	if e = state.Window.
		Size(20, 20).
		Title("glom, the visual code editor").
		Open().
		Run(func(gtx l.Context) l.Dimensions {
			return l.Dimensions{}
		},
			nil, func() {
				interrupt.Request()
			}, quit,
		); E.Chk(e) {
		
	}
}
