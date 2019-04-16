package main

import (
	"os"

	"github.com/crossXCheckIn/app"
	"github.com/crossXCheckIn/config"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		config.Logger.Error(err)
		os.Exit(0)
	}
	app.MakeCheckInForEverybody(config)
}
 