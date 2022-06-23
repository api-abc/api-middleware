package main

import (
	"log"
	"os"

	"github.com/api-abc/api-middleware/app"
	"github.com/api-abc/api-middleware/configuration"
)

func main() {
	var config configuration.Config
	if err := configuration.LoadConfig("config.json", &config); err != nil {
		log.Printf("unable to load config.json: %s", err)
		os.Exit(1)
	}

	di := configuration.NewDI(config)
	app.Run(di)

}
