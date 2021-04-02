package gel

import (
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/opts/binary"
	"github.com/p9c/qu"
)

type Theme struct {
	quit       qu.C
	shaper     text.Shaper
	collection Collection
	TextSize   unit.Value
	*Colors
	icons         map[string]*Icon
	scrollBarSize int
	Dark          *binary.Opt
	iconCache     IconCache
	WidgetPool    *Pool
}

// NewTheme creates a new theme to use for rendering a user interface
func NewTheme(dark *binary.Opt, fontCollection []text.FontFace, quit qu.C) (th *Theme) {
	th = &Theme{
		quit:          quit,
		shaper:        text.NewCache(fontCollection),
		collection:    fontCollection,
		TextSize:      unit.Sp(16),
		Colors:        newColors(),
		scrollBarSize: 0,
		iconCache:     make(IconCache),
	}
	th.SetDarkTheme(dark.True())
	return
}
