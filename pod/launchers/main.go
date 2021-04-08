package launchers

import (
	"crypto/tls"
	"fmt"
	"github.com/p9c/monorepo/node/state"
	"github.com/p9c/monorepo/pkg/apputil"
	"github.com/p9c/monorepo/pkg/chaincfg"
	"github.com/p9c/monorepo/pkg/chainrpc"
	"github.com/p9c/monorepo/pkg/fork"
	"github.com/p9c/monorepo/pkg/log"
	"github.com/p9c/monorepo/pkg/pipe/serve"
	"github.com/p9c/monorepo/pkg/pod"
	"github.com/p9c/monorepo/pkg/podopts"
	"github.com/p9c/monorepo/pkg/qu"
	"github.com/p9c/monorepo/pkg/util"
	"go.uber.org/atomic"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// // Main is the entrypoint for the pod suite
// func Main() int {
// 	log.SetLogLevel("trace")
// 	var e error
// 	var cx *pod.State
// 	if cx, e = GetNewContext(podcfgs.GetDefaultConfig()); F.Chk(e) {
// 		return 1
// 	}
//
// 	T.Ln("running command", cx.
// 		Config.
// 		RunningCommand.
// 		Name,
// 	)
// 	if e = cx.Config.RunningCommand.Entrypoint(cx); E.Chk(e) {
// 		return 1
// 	}
// 	return 0
// }

// GetNewContext returns a fresh new context
func GetNewContext(config *podopts.Config) (s *pod.State, e error) {
	// after this, all the configurations are set and mostly sanitized
	if e = config.Initialize(); E.Chk(e) {
		// return
		panic(e)
	}
	chainClientReady := qu.T()
	rand.Seed(time.Now().UnixNano())
	rand.Seed(rand.Int63())
	s = &pod.State{
		ChainClientReady: chainClientReady,
		KillAll:          qu.T(),
		Config:           config,
		ConfigMap:        config.Map,
		StateCfg:         new(state.Config),
		NodeChan:         make(chan *chainrpc.Server),
		Syncing:          atomic.NewBool(true),
	}
	if s.Config.RunningCommand.Colorizer != nil {
		log.AppColorizer = s.Config.RunningCommand.Colorizer
		log.App = s.Config.RunningCommand.AppText
	}
	// everything in the configuration is set correctly up to this point, except for settings based on the running
	// network, so after this is when those settings are elaborated
	I.Ln("setting active network:", s.Config.Network.V())
	switch s.Config.Network.V() {
	case "testnet", "testnet3", "t":
		s.ActiveNet = &chaincfg.TestNet3Params
		fork.IsTestnet = true
		// fork.HashReps = 3
	case "regtestnet", "regressiontest", "r":
		fork.IsTestnet = true
		s.ActiveNet = &chaincfg.RegressionTestParams
	case "simnet", "s":
		fork.IsTestnet = true
		s.ActiveNet = &chaincfg.SimNetParams
	default:
		if s.Config.Network.V() != "mainnet" &&
			s.Config.Network.V() != "m" {
			D.Ln("using mainnet for node")
		}
		s.ActiveNet = &chaincfg.MainNetParams
	}
	// if pipe logging is enabled, start it up
	if s.Config.PipeLog.True() {
		D.Ln("starting up pipe logger")
		serve.Log(s.KillAll, fmt.Sprint(os.Args))
	}
	// set to write logs in the network specific directory, if the value was set and is not the same as datadir
	if s.Config.LogDir.V() == s.Config.DataDir.V() {
		e = s.Config.LogDir.Set(filepath.Join(s.Config.DataDir.V(), s.ActiveNet.Name))
	}
	// set up TLS stuff if it hasn't been set up yet. We assume if the configured values correspond to files the files
	// are valid TLS cert/pairs, and that the key will be absent if onetimetlskey was set
	if (s.Config.ClientTLS.True() || s.Config.ServerTLS.True()) &&
		(
			(!apputil.FileExists(s.Config.RPCKey.V()) && s.Config.OneTimeTLSKey.False()) ||
				!apputil.FileExists(s.Config.RPCCert.V()) ||
				!apputil.FileExists(s.Config.CAFile.V())) {
		D.Ln("generating TLS certificates")
		I.Ln(s.Config.RPCKey.V(), s.Config.RPCCert.V(), s.Config.RPCKey.V())
		// Create directories for cert and key files if they do not yet exist.
		certDir, _ := filepath.Split(s.Config.RPCCert.V())
		keyDir, _ := filepath.Split(s.Config.RPCKey.V())
		e = os.MkdirAll(certDir, 0700)
		if e != nil {
			E.Ln(e)
			return
		}
		e = os.MkdirAll(keyDir, 0700)
		if e != nil {
			E.Ln(e)
			return
		}
		// Generate cert pair.
		org := "pod/wallet autogenerated cert"
		validUntil := time.Now().Add(time.Hour * 24 * 365 * 10)
		var cert, key []byte
		cert, key, e = util.NewTLSCertPair(org, validUntil, nil)
		if e != nil {
			E.Ln(e)
			return
		}
		_, e = tls.X509KeyPair(cert, key)
		if e != nil {
			E.Ln(e)
			return
		}
		// Write cert and (potentially) the key files.
		e = ioutil.WriteFile(s.Config.RPCCert.V(), cert, 0600)
		if e != nil {
			rmErr := os.Remove(s.Config.RPCCert.V())
			if rmErr != nil {
				E.Ln("cannot remove written certificates:", rmErr)
			}
			return
		}
		e = ioutil.WriteFile(s.Config.CAFile.V(), cert, 0600)
		if e != nil {
			rmErr := os.Remove(s.Config.RPCCert.V())
			if rmErr != nil {
				E.Ln("cannot remove written certificates:", rmErr)
			}
			return
		}
		e = ioutil.WriteFile(s.Config.RPCKey.V(), key, 0600)
		if e != nil {
			E.Ln(e)
			rmErr := os.Remove(s.Config.RPCCert.V())
			if rmErr != nil {
				E.Ln("cannot remove written certificates:", rmErr)
			}
			rmErr = os.Remove(s.Config.CAFile.V())
			if rmErr != nil {
				E.Ln("cannot remove written certificates:", rmErr)
			}
			return
		}
		D.Ln("done generating TLS certificates")
	}
	
	// Validate profile port number
	T.Ln("validating profile port number")
	if s.Config.Profile.V() != "" {
		var profilePort int
		profilePort, e = strconv.Atoi(s.Config.Profile.V())
		if e != nil || profilePort < 1024 || profilePort > 65535 {
			e = fmt.Errorf("the profile port must be between 1024 and 65535, disabling profiling")
			E.Ln(e)
			if e = s.Config.Profile.Set(""); E.Chk(e) {
			}
		}
	}
	
	T.Ln("checking addpeer and connectpeer lists")
	if s.Config.AddPeers.Len() > 0 && s.Config.ConnectPeers.Len() > 0 {
		e := fmt.Errorf("the addpeers and connectpeers options can not be mixed")
		_, _ = fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	
	T.Ln("checking proxy/connect for disabling listening")
	if (s.Config.ProxyAddress.V() != "" ||
		s.Config.ConnectPeers.Len() > 0) &&
		s.Config.P2PListeners.Len() == 0 {
		s.Config.DisableListen.T()
	}
	
	T.Ln("checking relay/reject nonstandard policy settings")
	switch {
	case s.Config.RelayNonStd.True() && s.Config.RejectNonStd.True():
		errf := "rejectnonstd and relaynonstd cannot be used together" +
			" -- choose only one, leaving neither activated"
		E.Ln(errf)
		// just leave both false
		s.Config.RelayNonStd.F()
		s.Config.RejectNonStd.F()
	case s.Config.RejectNonStd.True():
		s.Config.RelayNonStd.F()
	case s.Config.RelayNonStd.True():
		s.Config.RejectNonStd.F()
	}
	
	// Chk to make sure limited and admin users don't have the same username
	T.Ln("checking admin and limited username is different")
	if !s.Config.Username.Empty() &&
		s.Config.Username.V() == s.Config.LimitUser.V() {
		e := fmt.Errorf("--username and --limituser must not specify the same username")
		_, _ = fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	// Chk to make sure limited and admin users don't have the same password
	T.Ln("checking limited and admin passwords are not the same")
	if !s.Config.Password.Empty() &&
		s.Config.Password.V() == s.Config.LimitPass.V() {
		e := fmt.Errorf("password and limitpass must not specify the same password")
		_, _ = fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	
	
	return
}
