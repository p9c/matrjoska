package main

import (
	_ "gioui.org/app/permission/networkstate" // todo: integrate this into routeable package
	_ "gioui.org/app/permission/storage"      // this enables the home folder appdata directory to work on android (and ios)
	"github.com/p9c/matrjoska/pod/state"

	// _ "gioui.org/app/permission/bluetooth"
	// _ "gioui.org/app/permission/camera"
)

func main() {
	state.Main()
}
