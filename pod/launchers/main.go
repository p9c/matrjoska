package launchers

import (
	"github.com/p9c/monorepo/node/state"
	"github.com/p9c/monorepo/pkg/chainrpc"
	"github.com/p9c/monorepo/pkg/log"
	"github.com/p9c/monorepo/pkg/pod"
	"github.com/p9c/monorepo/pkg/podopts"
	"github.com/p9c/monorepo/pkg/qu"
	"github.com/p9c/monorepo/pod/podcfgs"
	"go.uber.org/atomic"
	"math/rand"
	"time"
)

// Main is the entrypoint for the pod suite
func Main() int {
	log.SetLogLevel("trace")
	var e error
	var cx *pod.State
	if cx, e = GetNewContext(podcfgs.GetDefaultConfig()); F.Chk(e) {
		return 1
	}

	T.Ln("running command", cx.
		Config.
		RunningCommand.
		Name,
	)
	if e = cx.Config.RunningCommand.Entrypoint(cx); E.Chk(e) {
		return 1
	}
	return 0
}

// GetNewContext returns a fresh new context
func GetNewContext(config *podopts.Config) (s *pod.State, e error) {
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
