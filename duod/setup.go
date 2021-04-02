package main

import (
	"fmt"
	"github.com/p9c/duod/pkg/chaincfg"
	"github.com/p9c/duod/pkg/chainrpc"
	"github.com/p9c/duod/pkg/fork"
	"github.com/p9c/duod/pkg/opts"
	"github.com/p9c/duod/pkg/pod"
	"github.com/p9c/duod/pkg/spec"
	"github.com/p9c/duod/pkg/state"
	"github.com/p9c/qu"
	"go.uber.org/atomic"
	"math/rand"
	"reflect"
	"time"
	
	// This enables pprof
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/trace"
	
	"github.com/p9c/duod/pkg/interrupt"
	"github.com/p9c/duod/pkg/limits"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 3)
	debug.SetGCPercent(10)
	var e error
	if runtime.GOOS != "darwin" {
		if e = limits.SetLimits(); E.Chk(e) { // todo: doesn't work on non-linux
			_, _ = fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", e)
			os.Exit(1)
		}
	}
	var f *os.File
	if os.Getenv("POD_TRACE") == "on" {
		D.Ln("starting trace")
		tpath := fmt.Sprintf("%v.trace", fmt.Sprint(os.Args))
		if f, e = os.Create(tpath); E.Chk(e) {
			E.Ln(
				"tracing env POD_TRACE=on but we can't write to trace file",
				fmt.Sprintf("'%s'", tpath),
				e,
			)
		} else {
			if e = trace.Start(f); E.Chk(e) {
			} else {
				D.Ln("tracing started")
				defer trace.Stop()
				defer func() {
					if e := f.Close(); E.Chk(e) {
					}
				}()
				interrupt.AddHandler(
					func() {
						D.Ln("stopping trace")
						trace.Stop()
						e := f.Close()
						if e != nil {
						}
					},
				)
			}
		}
	}
	var cx *pod.State
	if cx, e = GetNewContext(GetDefaultConfig()); E.Chk(e) {
		panic(e)
	}
	switch cx.Config.Network.V() {
	case "testnet", "testnet3", "t":
		cx.ActiveNet = &chaincfg.TestNet3Params
		fork.IsTestnet = true
		// fork.HashReps = 3
		cx.Config.Network.Set("testnet")
	case "regtestnet", "regressiontest", "r":
		fork.IsTestnet = true
		cx.ActiveNet = &chaincfg.RegressionTestParams
		cx.Config.Network.Set("regtestnet")
	case "simnet", "s":
		fork.IsTestnet = true
		cx.ActiveNet = &chaincfg.SimNetParams
		cx.Config.Network.Set("simnet")
	default:
		if cx.Config.Network.V() != "mainnet" &&
			cx.Config.Network.V() != "m" {
			D.Ln("using mainnet for node")
		}
		cx.Config.Network.Set("mainnet")
		cx.ActiveNet = &chaincfg.MainNetParams
	}
	
	cx.Config.Network.V()
	e = Main(cx)
	D.Ln("returning value", e, os.Args)
	if os.Getenv("POD_TRACE") == "on" {
		D.Ln("stopping trace")
		trace.Stop()
		// defer func() {
		// 	if e := f.Close(); E.Chk(e) {
		// 	}
		// }()
	}
	if E.Chk(e) {
		E.Ln("quitting with error")
		// D.Ln(interrupt.GoroutineDump())
		interrupt.Request()
		os.Exit(-1)
	}
}

// GetDefaultConfig returns a Config struct pristine factory fresh
func GetDefaultConfig() (c *opts.Config) {
	c = &opts.Config{
		Commands: spec.GetCommands(),
		Map:      spec.GetConfigs(),
	}
	c.RunningCommand = c.Commands[0]
	// I.S(c.Commands[0])
	// I.S(c.Map)
	t := reflect.ValueOf(c)
	t = t.Elem()
	for i := range c.Map {
		tf := t.FieldByName(i)
		if tf.IsValid() && tf.CanSet() && tf.CanAddr() {
			val := reflect.ValueOf(c.Map[i])
			tf.Set(val)
		}
	}
	// I.S(c)
	return
}

// GetNewContext returns a fresh new context
func GetNewContext(config *opts.Config) (s *pod.State, e error) {
	if e = config.Initialize(); E.Chk(e) {
		// return
		panic(e)
	}
	config.RunningCommand = config.Commands[0]
	chainClientReady := qu.T()
	rand.Seed(time.Now().UnixNano())
	rand.Seed(rand.Int63())
	s = &pod.State{
		ChainClientReady: chainClientReady,
		KillAll:          qu.T(),
		// App:              cli.NewApp(),
		Config:    config,
		ConfigMap: config.Map,
		StateCfg:  new(state.Config),
		// Language:         lang.ExportLanguage(appLang),
		// DataDir:          appdata.Dir(appName, false),
		NodeChan: make(chan *chainrpc.Server),
		Syncing:  atomic.NewBool(true),
	}
	return
}
