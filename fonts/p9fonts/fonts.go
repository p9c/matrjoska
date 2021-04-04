package p9fonts

import (
	"fmt"
	"sync"

	"gioui.org/font/opentype"
	"gioui.org/text"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	
	"github.com/p9c/monorepo/fonts/bariolbold"
	"github.com/p9c/monorepo/fonts/bariolbolditalic"
	"github.com/p9c/monorepo/fonts/bariollight"
	"github.com/p9c/monorepo/fonts/bariollightitalic"
	"github.com/p9c/monorepo/fonts/bariolregular"
	"github.com/p9c/monorepo/fonts/bariolregularitalic"
	"github.com/p9c/monorepo/fonts/plan9"
)

var (
	once       sync.Once
	collection []text.FontFace
	Fonts      = map[string]text.Font{
		"plan9":              {Typeface: "plan9"},
		"bariol regular":     {Typeface: "bariol regular"},
		"bariol italic":      {Typeface: "bariol italic", Style: text.Italic},
		"bariol bold":        {Typeface: "bariol bold", Weight: text.Bold},
		"bariol bolditalic":  {Typeface: "bariol bolditalic", Style: text.Italic, Weight: text.Bold},
		"bariol light":       {Typeface: "bariol light", Weight: text.Medium},
		"bariol lightitalic": {Typeface: "bariol lightitalic", Weight: text.Medium, Style: text.Italic},
		"go regular":         {Typeface: "go regular"},
		"go bold":            {Typeface: "go bold", Weight: text.Bold},
		"go bolditalic":      {Typeface: "go bolditalic", Weight: text.Bold, Style: text.Italic},
		"go italic":          {Typeface: "go italic", Style: text.Italic},
	}
)

func Collection() []text.FontFace {
	once.Do(func() {
		register(Fonts["plan9"], plan9.TTF)
		register(Fonts["bariol regular"], bariolregular.TTF)
		register(Fonts["bariol italic"], bariolregularitalic.TTF)
		register(Fonts["bariol bold"], bariolbold.TTF)
		register(Fonts["bariol bold italic"], bariolbolditalic.TTF)
		register(Fonts["bariol light"], bariollight.TTF)
		register(Fonts["bariol light italic"], bariollightitalic.TTF)
		register(Fonts["go regular"], gomono.TTF)
		register(Fonts["go bold"], gomonobold.TTF)
		register(Fonts["go bolditalic"], gomonobolditalic.TTF)
		register(Fonts["go italic"], gomonoitalic.TTF)
		// Ensure that any outside appends will not reuse the backing store.
		n := len(collection)
		collection = collection[:n:n]
	})
	return collection
}

func register(fnt text.Font, ttf []byte) {
	face, e := opentype.Parse(ttf)
	if e != nil  {
		panic(fmt.Errorf("failed to parse font: %v", e))
	}
	// fnt.Typeface = "Go"
	collection = append(collection, text.FontFace{Font: fnt, Face: face})
}