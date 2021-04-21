// +build !nominers

package launchers

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/gookit/color"

	"github.com/p9c/log"
	"github.com/p9c/matrjoska/cmd/kopach"
	"github.com/p9c/matrjoska/pod/state"

	"github.com/p9c/matrjoska/cmd/kopach/worker"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/fork"
	"github.com/p9c/matrjoska/pkg/interrupt"
)

// kopachHandle runs the kopach miner
func kopachHandle(ifc interface{}) (e error) {
	var cx *state.State
	var ok bool
	if cx, ok = ifc.(*state.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	// log.AppColorizer = color.Bit24(255, 128, 128, false).Sprint
	// log.App = "kopach"
	I.Ln("starting up kopach standalone miner for parallelcoin")
	D.Ln(os.Args)
	// podconfig.Configure(cx, true)
	if cx.ActiveNet.Name == chaincfg.TestNet3Params.Name {
		fork.IsTestnet = true
	}
	defer cx.KillAll.Q()
	e = kopach.Run(cx)
	<-interrupt.HandlersDone
	D.Ln("kopach main finished")
	return
}

func kopachWorkerHandle(cx *state.State) (e error) {
	log.AppColorizer = color.Bit24(255, 128, 128, false).Sprint
	log.App = "worker"
	if len(os.Args) > 3 {
		if os.Args[3] == chaincfg.TestNet3Params.Name {
			fork.IsTestnet = true
		}
	}
	if len(os.Args) > 4 {
		log.SetLogLevel(os.Args[4])
	}
	D.Ln("miner worker starting")
	w, conn := worker.New(os.Args[2], cx.KillAll, uint64(cx.Config.UUID.V()))
	e = rpc.Register(w)
	if e != nil {
		D.Ln(e)
		return e
	}
	D.Ln("starting up worker IPC")
	rpc.ServeConn(conn)
	D.Ln("stopping worker IPC")
	D.Ln("finished")
	return nil
}
