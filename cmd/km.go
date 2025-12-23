package main

import (
	"log"

	"github.com/kiemlicz/kmux/internal/common"
)

func main() {
	config, err := common.SetupConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}
	common.SetupLog(config.Log.Level)

}
