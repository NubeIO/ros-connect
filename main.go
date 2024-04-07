package main

import (
	"fmt"
	"github.com/NubeIO/ros-connect/backend"
	"github.com/NubeIO/ros-connect/backend/config"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	r := gin.Default()
	root, err := os.Getwd()
	dbPath := fmt.Sprintf("%s/db/db.yaml", root)
	configPath := fmt.Sprintf("%s/config/config.yaml", root)
	configInstance, err := config.New(configPath)
	opts := &backend.Options{
		Config: configInstance,
	}
	instanceManager, err := backend.New(dbPath, configPath, opts)
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		log.Fatalf("init appStore on start of app err: %s", err.Error())
	}
	// Setup routes
	instanceManager.SetupRoutes(r)
	r.Run(":1771")
}
