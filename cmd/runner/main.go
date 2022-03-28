package main

import (
	"flag"
	"log"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/logger"
	"github.com/utkuufuk/entrello/internal/service"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "config.json", "config file path")
	flag.Parse()

	cfg, err := config.ReadRunnerConfig(configFile)
	if err != nil {
		log.Fatalf("Could not read configuration: %v", err)
	}

	if err = service.Poll(cfg); err != nil {
		logger.Error(err.Error())
	}
}
