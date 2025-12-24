package main

import (
	"log"
	"os"

	"github.com/kiemlicz/kmux/internal/common"
	"github.com/kiemlicz/kmux/internal/kmux"
)

func main() {
	config, ops, err := common.SetupConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}
	common.SetupLog(config.Log.Level)
	common.Log.Debugf("Config: %+v", config)

	km := kmux.NewKmux(*config)

	if ops.Start != "" {
		err := km.StartEnvironment(*ops)
		if err != nil {
			common.Log.Errorf("Failed to start environment: %v", err)
			os.Exit(8)
		}
	} else if ops.Discover != "" {
		err := km.DiscoverEnvironment(*ops)
		if err != nil {
			common.Log.Errorf("Failed to discover environment namespaces: %v", err)
			os.Exit(9)
		}
	} else if ops.New != "" {
		err := km.NewEnvironment(ops)
		if err != nil {
			common.Log.Errorf("Failed to create new environment: %v", err)
			os.Exit(10)
		}
		common.Log.Infof("Environment created, start and populate KUBECONFIG (%s)", ops.Kubeconfig)
	} else {
		common.Log.Error("No supported command provided")
		os.Exit(11)
	}
}
