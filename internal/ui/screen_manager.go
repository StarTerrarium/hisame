package ui

import (
	"fyne.io/fyne/v2"
	"github.com/StarTerrarium/hisame/internal/state"
	"github.com/sirupsen/logrus"
	"sync"
)

// ScreenManager acts as a central management tool for changing between the main screens available in the app.
type ScreenManager struct {
	window     fyne.Window
	mainScreen *MainScreen
	isAuth     bool
}

var (
	instance          *ScreenManager
	screenManagerOnce sync.Once
)

// GetScreenManager returns the singleton instance of ScreenManager.
// It panics if InitialiseScreenManager has not been called yet.
func getScreenManager() *ScreenManager {
	if instance == nil {
		logrus.Panic("Attempt to access ScreenManager before initialising.  This should never happen and is a bug in the code.  Please open an issue.")
		panic("Attempt to access ScreenManager before initialising.")
	}
	return instance
}

func InitialiseScreenManager(window fyne.Window) {
	screenManagerOnce.Do(func() {
		instance = &ScreenManager{
			window: window,
			isAuth: state.GetAppState().GetAuthToken() != "",
		}
		instance.mainScreen = NewMainScreen(window)
		instance.showInitialPage()
	})
}

func (sm *ScreenManager) showInitialPage() {
	if sm.isAuth {
		sm.ShowPage(NewAnimeListPage())
	} else {
		sm.ShowPage(NewLoginPage())
	}
}

func (sm *ScreenManager) ShowPage(page Page) {
	sm.mainScreen.ShowPage(page)
}

func (sm *ScreenManager) HandleLoginSuccess(token string) {
	sm.isAuth = true
	appState := state.GetAppState()
	appState.SetAuthToken(token)
	err := appState.SaveAuthToken()
	if err != nil {
		logrus.Errorf("Failed to save auth token: %v", err)
	}
	// Enable navigation buttons
	sm.mainScreen.navigationBar.UpdateAuthenticationState(sm.isAuth)
	sm.ShowPage(NewAnimeListPage())
}

func (sm *ScreenManager) HandleLogout() {
	sm.isAuth = false
	// Clear authentication tokens and state
	err := state.GetAppState().ClearAuthToken()
	if err != nil {
		logrus.Errorf("Failed to clear auth token: %v", err)
	}
	// Disable navigation buttons
	sm.mainScreen.navigationBar.UpdateAuthenticationState(sm.isAuth)
	sm.ShowPage(NewLoginPage())
}
