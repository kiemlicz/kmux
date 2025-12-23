package main

import (
	"log"
	"os"

	"github.com/kiemlicz/kmux/internal/common"
	"github.com/kiemlicz/kmux/internal/kmux"
)

func main() {
	config, err := common.SetupConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}
	common.SetupLog(config.Log.Level)
	common.Log.Debugf("Config: %+v", config)

	if config.Start != "" {
		err := kmux.StartEnvironment(*config)
		if err != nil {
			common.Log.Errorf("Failed to start environment: %v", err)
			os.Exit(8)
		}
	} else if config.Discover != "" {
		err := kmux.DiscoverEnvironment(*config)
		if err != nil {
			common.Log.Errorf("Failed to discover environment namespaces: %v", err)
			os.Exit(9)
		}
	} else if config.New != "" {
		err := kmux.NewEnvironment(config)
		if err != nil {
			common.Log.Errorf("Failed to create new environment: %v", err)
			os.Exit(10)
		}
		err = common.SaveConfig(config)
		if err != nil {
			common.Log.Errorf("Failed to save configuration: %v", err)
			os.Exit(8)
		}
		common.Log.Infof("Environment created, start and populate KUBECONFIG (%s)", config.Kubeconfig)
	} else {
		common.Log.Error("No supported command provided")
		os.Exit(11)
	}
}
