package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/StarTerrarium/hisame/internal/auth"
	"github.com/StarTerrarium/hisame/internal/config"
	"github.com/StarTerrarium/hisame/internal/state"
	"github.com/StarTerrarium/hisame/internal/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	// TODO:  Make log level configurable from config file
	cleanupLogger := utils.InitLogger()
	defer cleanupLogger()

	// Load user config from file
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Warnf("Error loading config.  Will use default config. %v", err)
		cfg = config.DefaultConfig()
	}

	appState := state.InitialiseAppState(cfg)
	// Temporary, just to get this to compile.
	appState.GetConfig()

	logrus.Infof("App state initialised.  Log level: %s", logrus.GetLevel().String())

	a := app.New()
	w := a.NewWindow("Hisame")

	// Load token if exists & check if expiring soon
	authenticated := false // Placeholder for now, force login every time

	if !authenticated {
		err := auth.Login()
		if err != nil {
			logrus.Fatalf("failed to login: %s", err)
			return
		}
	}

	w.ShowAndRun()
}
