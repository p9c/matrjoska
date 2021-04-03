package pathtree

import (
	l "gioui.org/layout"
	"github.com/p9c/monorepo/gel"
)

type Widget struct {
	*gel.Window
}

func New(w *gel.Window) *Widget {
	return &Widget{Window: w}
}

func (w *Widget) Fn(gtx l.Context) l.Dimensions {
	return w.H3("glom").Fn(gtx)
}
