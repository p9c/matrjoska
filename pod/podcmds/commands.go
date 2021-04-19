package podcmds

import (
	"fmt"

	"github.com/gookit/color"

	"github.com/p9c/matrjoska/pod/launchers"
	"github.com/p9c/matrjoska/version"
	"github.com/p9c/opts/cmds"
)

// GetCommands returns available subcommands in Parallelcoin Pod
func GetCommands() (c cmds.Commands) {
	c = cmds.Commands{
		{Name: "gui", Description:
		"ParallelCoin GUI Wallet/Miner/Explorer",
			Entrypoint: func(c interface{}) error { return nil },
			Colorizer:  color.Bit24(128, 255, 255, false).Sprint,
			AppText:    "   gui",
		},
		{Name: "version", Description:
		"print version and exit",
			Entrypoint: func(c interface{}) error {
				fmt.Println(version.Tag)
				return nil
			},
		},
		{Name: "ctl", Description:
		"command line wallet and chain RPC client",
			Entrypoint: launchers.CtlHandle,
			Colorizer:  color.Bit24(128, 255, 128, false).Sprint,
			AppText:    "   ctl",
			Commands: []cmds.Command{
				{Name: "list", Description:
				"list available commands",
					Entrypoint: launchers.CtlHandleList,
				},
			},
		},
		{Name: "node", Description:
		"ParallelCoin blockchain node",
			Entrypoint: launchers.NodeHandle,
			Colorizer:  color.Bit24(128, 128, 255, false).Sprint,
			AppText:    "  node",
			Commands: []cmds.Command{
				{Name: "dropaddrindex", Description:
				"drop the address database index",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "droptxindex", Description:
				"drop the transaction database index",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "dropcfindex", Description:
				"drop the cfilter database index",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "dropindexes", Description:
				"drop all of the indexes",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "resetchain", Description:
				"deletes the current blockchain cache to force redownload",
					Entrypoint: func(c interface{}) error { return nil },
				},
			},
		},
		{Name: "wallet", Description:
		"run the wallet server (requires a chain node to function)",
			Entrypoint: launchers.WalletHandle, // func(c interface{}) error { return nil },
			Commands: []cmds.Command{
				{Name: "drophistory", Description:
				"reset the wallet transaction history",
					Entrypoint: func(c interface{}) error { return nil },
				},
			},
			Colorizer: color.Bit24(255, 255, 128, false).Sprint,
			AppText:   "wallet",
		},
		{Name: "kopach", Description:
		"standalone multicast miner for easy mining farm deployment",
			Entrypoint: func(c interface{}) error { return nil },
			Colorizer:  color.Bit24(255, 128, 128, false).Sprint,
			AppText:    "kopach",
		},
		{Name: "worker", Description:
		"single thread worker process, normally started by kopach",
			Entrypoint: func(c interface{}) error { return nil },
			Colorizer:  color.Bit24(255, 128, 255, false).Sprint,
			AppText:    "worker",
		},
	}
	c.PopulateParents(nil)
	// I.S(c)
	return
}

