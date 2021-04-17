package launchers

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/p9c/matrjoska/ctl"
	"github.com/p9c/matrjoska/kopach"
	"github.com/p9c/matrjoska/kopach/worker"
	"github.com/p9c/matrjoska/node/node"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/constant"
	"github.com/p9c/matrjoska/pkg/fork"
	"github.com/p9c/matrjoska/pkg/interrupt"
	"github.com/p9c/matrjoska/pkg/log"
	"github.com/p9c/matrjoska/pkg/pod"
	"io/ioutil"
	"net/rpc"
	"os"
	"path/filepath"
	
	"github.com/p9c/matrjoska/pkg/qu"
	
	"github.com/p9c/matrjoska/pkg/apputil"
	"github.com/p9c/matrjoska/walletmain"
)

// NodeHandle runs the ParallelCoin blockchain node
func NodeHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	// log.AppColorizer = color.Bit24(128, 128, 255, false).Sprint
	// log.App = "  node"
	I.Ln("running node handler")
	// podconfig.Configure(cx, true)
	cx.NodeReady = qu.T()
	cx.Node.Store(false)
	// // serviceOptions defines the configuration options for the daemon as a service on Windows.
	// type serviceOptions struct {
	// 	ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
	// }
	// // runServiceCommand is only set to a real function on Windows. It is used to parse and execute service commands
	// // specified via the -s flag.
	// runServiceCommand := func(string) (e error) { return nil }
	// // Service options which are only added on Windows.
	// serviceOpts := serviceOptions{}
	// // Perform service command and exit if specified. Invalid service commands show an appropriate error. Only runs
	// // on Windows since the runServiceCommand function will be nil when not on Windows.
	// if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
	// 	if e = runServiceCommand(serviceOpts.ServiceCommand); E.Chk(e) {
	// 		return e
	// 	}
	// 	return nil
	// }
	// config.Configure(cx, c.Command.Name, true)
	// D.Ln("starting shell")
	// if cx.Config.ClientTLS.True() || cx.Config.ServerTLS.True() {
	// 	// generate the tls certificate if configured
	// 	if apputil.FileExists(cx.Config.RPCCert.V()) &&
	// 		apputil.FileExists(cx.Config.RPCKey.V()) &&
	// 		apputil.FileExists(cx.Config.CAFile.V()) {
	// 	} else {
	// 		if _, e = walletmain.GenerateRPCKeyPair(cx.Config, true); E.Chk(e) {
	// 		}
	// 	}
	// }
	// if cx.Config.NodeOff.False() {
	go func() {
		if e := node.Main(cx); E.Chk(e) {
			E.Ln("error starting node ", e)
		}
	}()
	I.Ln("starting node")
	if cx.Config.DisableRPC.False() {
		cx.RPCServer = <-cx.NodeChan
		cx.NodeReady.Q()
		cx.Node.Store(true)
		I.Ln("node started")
	}
	// }
	cx.WaitWait()
	I.Ln("node is now fully shut down")
	cx.WaitGroup.Wait()
	<-cx.KillAll
	return nil
}

// walletHandle runs the wallet server
func walletHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	// log.AppColorizer = color.Bit24(255, 255, 128, false).Sprint
	// log.App = "wallet"
	// podconfig.Configure(cx, true)
	cx.Config.WalletFile.Set(filepath.Join(cx.Config.DataDir.V(), cx.ActiveNet.Name, constant.DbName))
	// dbFilename := *cx.Config.DataDir + slash + cx.ActiveNet.
	// 	Params.Name + slash + wallet.WalletDbName
	if !apputil.FileExists(cx.Config.WalletFile.V()) && !cx.IsGUI {
		// D.Ln(cx.ActiveNet.Name, *cx.Config.WalletFile)
		if e = walletmain.CreateWallet(cx.ActiveNet, cx.Config); E.Chk(e) {
			E.Ln("failed to create wallet", e)
			return e
		}
		fmt.Println("restart to complete initial setup")
		os.Exit(0)
	}
	// for security with apps launching the wallet, the public password can be set with a file that is deleted after
	walletPassPath := filepath.Join(cx.Config.DataDir.V(), cx.ActiveNet.Name, "wp.txt")
	D.Ln("reading password from", walletPassPath)
	if apputil.FileExists(walletPassPath) {
		var b []byte
		if b, e = ioutil.ReadFile(walletPassPath); !E.Chk(e) {
			cx.Config.WalletPass.SetBytes(b)
			D.Ln("read password '" + string(b) + "'")
			for i := range b {
				b[i] = 0
			}
			if e = ioutil.WriteFile(walletPassPath, b, 0700); E.Chk(e) {
			}
			if e = os.Remove(walletPassPath); E.Chk(e) {
			}
			D.Ln("wallet cookie deleted", *cx.Config.WalletPass)
		}
	}
	cx.WalletKill = qu.T()
	if e = walletmain.Main(cx); E.Chk(e) {
		E.Ln("failed to start up wallet", e)
	}
	// if !*cx.Config.DisableRPC {
	// 	cx.WalletServer = <-cx.WalletChan
	// }
	// cx.WaitGroup.Wait()
	cx.WaitWait()
	return
}

// kopachHandle runs the kopach miner
func kopachHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
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

func kopachWorkerHandle(cx *pod.State) (e error) {
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

func CtlHandleList(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	_ = cx
	// fmt.Println("Here are the available commands. Pausing a moment as it is a long list...")
	// time.Sleep(2 * time.Second)
	ctl.ListCommands()
	return nil
}

func CtlHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	// log.AppColorizer = color.Bit24(128, 128, 255, false).Sprint
	// log.App = "   ctl"
	cx.Config.LogLevel.Set("off")
	// podconfig.Configure(cx, true)
	ctl.Main(cx)
	return nil
}
