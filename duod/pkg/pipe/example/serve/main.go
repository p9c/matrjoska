package main

import (
	"fmt"
	"time"
	
	"github.com/p9c/monorepo/qu"
	
	"github.com/p9c/monorepo/monorepo/duod/pkg/pipe"
)

func main() {
	p := pipe.Serve(qu.T(), func(b []byte) (e error) {
		fmt.Print("from parent: ", string(b))
		return
	},
	)
	for {
		_, e := p.Write([]byte("ping"))
		if e != nil {
			fmt.Println("err:", e)
		}
		time.Sleep(time.Second)
	}
}
