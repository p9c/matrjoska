package old

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gookit/color"

	"github.com/p9c/log"
	"github.com/p9c/matrjoska/pkg/constant"
	"github.com/p9c/matrjoska/pkg/wallet"
	"github.com/p9c/matrjoska/pod/state"
	"github.com/p9c/opts"

	"github.com/urfave/cli"

	"github.com/p9c/matrjoska/pod/podconfig"

	"github.com/p9c/matrjoska/cmd/node"
	"github.com/p9c/matrjoska/pkg/apputil"
)

func ShellHandle(cx *state.State) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		log.AppColorizer = color.Bit24(255, 128, 128, false).Sprint
		log.App = " shell"
		podconfig.Configure(cx, true)
		opts.D.Ln("starting shell")
		if cx.Config.TLS.True() || cx.Config.ServerTLS.True() {
			// generate the tls certificate if configured
			if apputil.FileExists(cx.Config.RPCCert.V()) && apputil.FileExists(cx.Config.RPCKey.V()) &&
				apputil.FileExists(cx.Config.CAFile.V()) {
				
			} else {
				_, _ = wallet.GenerateRPCKeyPair(cx.Config, true)
			}
		}
		dbFilename := filepath.Join(
			cx.Config.DataDir.V(),
			cx.ActiveNet.Name,
			constant.DbName,
		)
		if !apputil.FileExists(dbFilename) && !cx.IsGUI {
			// log.SetLevel("off", false)
			if e := wallet.CreateWallet(cx.ActiveNet, cx.Config); opts.E.Chk(e) {
				opts.E.Ln("failed to create wallet", e)
			}
			fmt.Println("restart to complete initial setup")
			os.Exit(1)
		}
		// for security with apps launching the wallet, the public password can be set with a file that is deleted after
		walletPassPath := filepath.Join(cx.Config.DataDir.V(), cx.ActiveNet.Name, "wp.txt")
		opts.D.Ln("reading password from", walletPassPath)
		if apputil.FileExists(walletPassPath) {
			var b []byte
			if b, e = ioutil.ReadFile(walletPassPath); !opts.E.Chk(e) {
				cx.Config.WalletPass.SetBytes(b)
				opts.D.Ln("read password '" + string(b) + "'")
				for i := range b {
					b[i] = 0
				}
				if e = ioutil.WriteFile(walletPassPath, b, 0700); opts.E.Chk(e) {
				}
				if e = os.Remove(walletPassPath); opts.E.Chk(e) {
				}
				opts.D.Ln("wallet cookie deleted", *cx.Config.WalletPass)
			}
		}
		if cx.Config.NodeOff.False() {
			go func() {
				e = node.Main(cx)
				if e != nil {
					opts.E.Ln("error starting node ", e)
				}
			}()
			opts.I.Ln("starting node")
			if cx.Config.DisableRPC.False() {
				cx.RPCServer = <-cx.NodeChan
			}
			opts.I.Ln("node started")
		}
		if cx.Config.WalletOff.False() {
			go func() {
				e = wallet.Main(cx)
				if e != nil {
					fmt.Println("error running wallet:", e)
				}
			}()
			// I.Ln("starting wallet")
			// if !*cx.Config.DisableRPC {
			// 	cx.WalletServer = <-cx.WalletChan
			// }
			// I.Ln("wallet started")
		}
		opts.D.Ln("shell started")
		// cx.WaitGroup.Wait()
		cx.WaitWait()
		return nil
	}
}