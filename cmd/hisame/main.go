package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/StarTerrarium/hisame/internal/auth"
	"github.com/StarTerrarium/hisame/internal/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	// TODO:  Make log level configurable from config file
	cleanupLogger := utils.InitLogger()
	defer cleanupLogger()

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
