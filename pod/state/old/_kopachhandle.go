package old

import (
	"fmt"
	"os"

	"github.com/gookit/color"

	"github.com/p9c/log"
	"github.com/p9c/matrjoska/cmd/kopach"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/fork"
	"github.com/p9c/matrjoska/pod/state"
	"github.com/p9c/opts"

	"github.com/p9c/interrupt"

	"github.com/p9c/matrjoska/pod/podconfig"
)

// kopachHandle runs the kopach miner
func KopachHandle(ifc interface{}) (e error) {
	var cx *state.State
	var ok bool
	if cx, ok = ifc.(*state.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	log.AppColorizer = color.Bit24(255, 128, 128, false).Sprint
	log.App = "kopach"
	opts.I.Ln("starting up kopach standalone miner for parallelcoin")
	opts.D.Ln(os.Args)
	podconfig.Configure(cx, true)
	if cx.ActiveNet.Name == chaincfg.TestNet3Params.Name {
		fork.IsTestnet = true
	}
	defer cx.KillAll.Q()
	e = kopach.Run(cx)
	<-interrupt.HandlersDone
	opts.D.Ln("kopach main finished")
	return
}
