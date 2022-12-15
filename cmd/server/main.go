package main

import (
	"github.com/JeffersonQin/syncat/internal/server"
	"github.com/JeffersonQin/syncat/pkg/config"
	"github.com/JeffersonQin/syncat/pkg/database"
	"log"
)

func init() {
	// Load common config
	log.Println("Loading config...")
	err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed to load config.", err)
	}
	// Load server config
	log.Println("Loading server config...")
	err = server.LoadConfig()
	if err != nil {
		log.Fatalln("failed to load server config.", err)
	}
}

func main() {
	// Load database
	log.Println("Loading database...")
	err := database.LoadDatabase()
	if err != nil {
		log.Fatalln("failed to open database.", err)
	}
	defer func() {
		err := database.CloseDatabase()
		if err != nil {
			log.Println("failed to close database.", err)
		}
	}()

	// Start server
	log.Println("Starting server...")
	err = server.StartSyncatServer()
	if err != nil {
		log.Fatalln("failed to start syncat server.", err)
	}
}
