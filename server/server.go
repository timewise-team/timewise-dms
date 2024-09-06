package server

import (
	"dbms/config"
	"dbms/database"
	"log"
)

func RegisterServer() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	_, err = database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

}
