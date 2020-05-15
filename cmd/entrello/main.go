package main

import (
	"fmt"
	"log"

	"github.com/utkuufuk/entrello/internal/config"
)

func main() {
	cfg, err := config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("[-] could not read config variables: %v", err)
	}

	fmt.Printf("%v", cfg)
}
