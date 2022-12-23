package main

import (
	"fmt"
	"os"

	"github.com/MetalBlockchain/antelopevm/vm"
	"github.com/MetalBlockchain/metalgo/utils/logging"
	"github.com/MetalBlockchain/metalgo/utils/ulimit"
	"github.com/MetalBlockchain/metalgo/vms/rpcchainvm"
	log "github.com/inconshreveable/log15"
)

func main() {
	version, err := PrintVersion()

	if err != nil {
		fmt.Printf("couldn't get config: %s", err)
		os.Exit(1)
	}

	if version {
		fmt.Println(vm.Version)
		os.Exit(0)
	}

	if err := ulimit.Set(ulimit.DefaultFDLimit, logging.NoLog{}); err != nil {
		fmt.Printf("failed to set fd limit correctly due to: %s", err)
		os.Exit(1)
	}

	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat())))

	rpcchainvm.Serve(&vm.VM{})
}
