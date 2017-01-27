package main

import (
	"flag"
	"os"

	"./config/"
	"./webservice"
	log "github.com/Sirupsen/logrus"
)

var SystemConfig *config.SystemConfig

func main() {

	debugFlag := flag.Bool("debug", false, "enable debug log")

	flag.Parse()

	//log.SetFlags(log.LstdFlags | log.Lshortfile)

	//config.Log = log.New(os.Stderr, "config-ws: ", log.LstdFlags | log.Lshortfile)
	SystemConfig = config.ReadFile("etc/system.conf")

	SystemConfig.Log = log.New()
	SystemConfig.Log.Out = os.Stderr

	if *debugFlag {
		SystemConfig.Log.Level = log.DebugLevel
		SystemConfig.Log.Debug("Log level debug activated")

	} else {
		SystemConfig.Log.Info("Log level info")
	}

	webservice.Start(SystemConfig)

	//ReadJson("example.json")

	os.Exit(0)
}
