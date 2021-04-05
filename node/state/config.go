package state

import (
	"github.com/p9c/monorepo/pkg/amt"
	"github.com/p9c/monorepo/pkg/btcaddr"
	"github.com/p9c/monorepo/pkg/connmgr"
	"net"
	"time"
	
	"github.com/p9c/monorepo/pkg/chaincfg"
)

// Config stores current state of the node
type Config struct {
	Lookup              connmgr.LookupFunc
	Oniondial           func(string, string, time.Duration) (net.Conn, error)
	Dial                func(string, string, time.Duration) (net.Conn, error)
	AddedCheckpoints    []chaincfg.Checkpoint
	ActiveMiningAddrs   []btcaddr.Address
	ActiveMinerKey      []byte
	ActiveMinRelayTxFee amt.Amount
	ActiveWhitelists    []*net.IPNet
	DropAddrIndex       bool
	DropTxIndex         bool
	DropCfIndex         bool
	Save                bool
}
