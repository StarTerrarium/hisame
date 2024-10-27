package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/StarTerrarium/hisame/internal/config"
	"github.com/StarTerrarium/hisame/internal/state"
	"github.com/StarTerrarium/hisame/internal/ui"
	"github.com/StarTerrarium/hisame/internal/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	cleanupLogger := utils.InitLogger()
	defer cleanupLogger()

	// Load user config from file
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Warnf("Error loading config.  Will use default config. %v", err)
		cfg = config.DefaultConfig()
	}

	state.InitialiseAppState(cfg)

	logrus.Infof("App state initialised.  Log level: %s", logrus.GetLevel().String())

	a := app.NewWithID("Hisame")
	w := a.NewWindow("Hisame")
	// Starting with a huge size makes the application start maximised, at least in KDE plasma.
	// TODO: Confirm behaviour on other DE & OS
	w.Resize(fyne.NewSize(7680, 4320))

	ui.InitialiseScreenManager(w)

	logrus.Info("Starting GUI")
	w.ShowAndRun()
}
