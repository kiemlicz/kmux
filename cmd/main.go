package main

import (
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
		err := km.StartEnvironment(*ops)
		if err != nil {
			common.Log.Errorf("Failed to start environment: %v", err)
			os.Exit(8)
		}
	case common.OptionDiscover:
		err := km.DiscoverEnvironment(*ops)
		if err != nil {
			common.Log.Errorf("Failed to discover environment namespaces: %v", err)
			os.Exit(9)
		}
	case common.OptionNew:
		err := km.NewEnvironment(ops)
		if err != nil {
			common.Log.Errorf("Failed to create new environment: %v", err)
			os.Exit(10)
		}
		common.Log.Infof("Environment created, start and populate KUBECONFIG (%s)", ops.Kubeconfig)
	case common.OptionCompletions:
		completions, err := kmux.CompletionsZsh(config)
		if err != nil {
			common.Log.Errorf("Failed to generate completions: %v", err)
			os.Exit(12)
		}
		common.Log.Infof("Generated ZSH completions, paste to sourced zsh rc file:\n%s", completions)
	default:
		common.Log.Error("No supported command provided")
		os.Exit(11)
	}
}
