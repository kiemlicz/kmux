package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kiemlicz/kmux/internal/common"
	"github.com/kiemlicz/kmux/internal/kmux"
)

var Version = "dev"

func main() {
	config, ops, err := common.SetupConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}
	common.SetupLog(config.Log.Level)
	common.Log.Debugf("KMux: %s, Config: %+v, Ops: %+v", Version, config, ops)

	km := kmux.NewKmux(*config)

	switch ops.OperationName {
	case common.OptionStart:
		if err := km.StartEnvironment(*ops); err != nil {
			common.Log.Errorf("Failed to start environment: %v", err)
			os.Exit(8)
		}
	case common.OptionDiscover:
		if err := km.DiscoverEnvironment(*ops); err != nil {
			common.Log.Errorf("Failed to discover environment namespaces: %v", err)
			os.Exit(9)
		}
	case common.OptionNew:
		if err := km.NewEnvironment(ops); err != nil {
			common.Log.Errorf("Failed to create new environment: %v", err)
			os.Exit(10)
		}
		// reload environment to allow for discover or handle differently
		if err = km.DiscoverEnvironment(*ops); err != nil {
			common.Log.Errorf("Failed to run discover on new environment: %v, ensure KUBECONFIG: %s exists", err, ops.Kubeconfig)
			os.Exit(10)
		}
		common.Log.Infof("Environment created, KUBECONFIG (%s)", ops.Kubeconfig)
	case common.OptionCompletions:
		completions, err := kmux.CompletionsZsh(config)
		if err != nil {
			common.Log.Errorf("Failed to generate completions: %v", err)
			os.Exit(12)
		}
		fmt.Println(completions)
	default:
		common.Log.Error("No supported command provided")
		os.Exit(11)
	}
}
