// +build headless

package launchers

import (
	"os"

	"github.com/p9c/matrjoska/pod/state"
)

func GUIHandle(ifc interface{}) (e error) {
	state.W.Ln("GUI was disabled for this podbuild (server only version)")
	os.Exit(1)
	return nil
}
