package main

import (
	"github.com/JeffersonQin/syncat/pkg/config"
	"log"
)

func init() {
	// Load common config
	log.Println("Loading config...")
	err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed to load config.", err)
	}
}

func main() {

}
