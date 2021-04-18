// +build nominers

package launchers

import (
	"github.com/p9c/matrjoska/pkg/pod"
)

// kopachHandle runs the kopach miner
func kopachHandle(ifc interface{}) (e error) {
	D.Ln("kopach disabled for ios/android")
	return
}

func kopachWorkerHandle(cx *pod.State) (e error) {
	D.Ln("kopach worker disabled for ios/android")
	return nil
}
