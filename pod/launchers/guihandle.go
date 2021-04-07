// +build !headless

package launchers

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/p9c/monorepo/pkg/log"
	"github.com/p9c/monorepo/pkg/pod"
	
	"github.com/p9c/monorepo/gui"
	"github.com/p9c/monorepo/pod/podconfig"
)

// GUIHandle starts up the GUI wallet
func GUIHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	log.AppColorizer = color.Bit24(128, 255, 255, false).Sprint
	log.App = "   gui"
	D.Ln("starting up parallelcoin pod gui...")
	// fork.ForkCalc()
	podconfig.Configure(cx, true)
	// D.Ln(os.Args)
	// interrupt.AddHandler(func() {
	// 	D.Ln("wallet gui is shut down")
	// })
	if e = gui.Main(cx); E.Chk(e) {
	}
	D.Ln("pod gui finished")
	return
}