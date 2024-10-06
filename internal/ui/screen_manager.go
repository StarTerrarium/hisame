package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/StarTerrarium/hisame/internal/auth"
	"github.com/sirupsen/logrus"
)

// ScreenManager acts as a central management tool for changing between the main screens available in the app.
type ScreenManager struct {
	window fyne.Window
}

func NewScreenManager(w fyne.Window) *ScreenManager {
	return &ScreenManager{
		window: w,
	}
}

// ShowLoginScreen prompts the user to login.  Accepts an onsuccess callback function
// to execute when the login has succeeded.
func (sm *ScreenManager) ShowLoginScreen(onSuccess func()) {
	// TODO:  Move into its own file?  Or is this acceptably small..
	authInstance := auth.NewAuth()

	loginButton := widget.NewButton("Login with AniList", func() {
		logrus.Infof("Starting login.  Login URL: %s", authInstance.LoginURL)
		err := authInstance.StartCallbackServer()
		if err != nil {
			logrus.Errorf("Error starting login flow: %v", err)
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Login error",
				Content: "Error starting Login flow.  Please check logs and try again.",
			})
			// TODO:  Figure out how to do proper error feedback here.  Notification is good enough for now.
			return
		}

		// Needs to be referenced inside the cancel button callback, to hide the dialog
		var loadingDialog *widget.PopUp

		loadingLabel := widget.NewLabel("Logging in..  Continue in your browser")
		loadingSpinner := widget.NewProgressBarInfinite()
		manualLink := widget.NewHyperlink("If the browser didn't open, click here to login manually", authInstance.LoginURL)
		cancelButton := widget.NewButton("Cancel", func() {
			logrus.Info("Login cancelled by user")
			authInstance.StopCallbackServer()
			loadingDialog.Hide()
		})

		loadingContent := container.NewVBox(loadingLabel, loadingSpinner, manualLink, cancelButton)
		loadingDialog = widget.NewModalPopUp(loadingContent, sm.window.Canvas())
		loadingDialog.Show()

		err = fyne.CurrentApp().OpenURL(authInstance.LoginURL)
		if err != nil {
			logrus.Warnf("Error opening Login URL: %v", err)
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Error opening Login URL",
				Content: "Hisame was unable to open the AniList login page in your browser.",
			})
			// Do not return, we want to let the user try to manually action the login.  They can click cancel
			// on the modal to exit if desired.
		}

		go func() {
			token, err := authInstance.WaitForToken()
			if err != nil {
				logrus.Error("Error waiting for token", err)
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Login error",
					Content: "There was an error reading the auth token.  Please check the logs and try again.",
				})
				loadingDialog.Hide()
				return
			}
			logrus.Tracef("Received token: %s", token)

			loadingDialog.Hide()
			onSuccess()
		}()
	})

	loginContent := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(
			container.NewVBox(
				widget.NewLabel("Login to AniList to use Hisame"),
				loginButton,
			),
		),
		layout.NewSpacer(),
	)
	sm.window.SetContent(loginContent)
}
