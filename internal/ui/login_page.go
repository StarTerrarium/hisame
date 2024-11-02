package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/StarTerrarium/hisame/internal/auth"
	"github.com/sirupsen/logrus"
)

type LoginPage struct {
	content fyne.CanvasObject
}

func NewLoginPage() *LoginPage {
	lp := &LoginPage{}
	lp.content = lp.buildContent()
	return lp
}

func (lp *LoginPage) Content() fyne.CanvasObject {
	return lp.content
}

func (lp *LoginPage) buildContent() fyne.CanvasObject {
	authInstance := auth.NewAuth()

	loginButton := widget.NewButton("Login with AniList", func() {
		lp.startLoginFlow(authInstance)
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
	return loginContent
}

func (lp *LoginPage) startLoginFlow(authInstance *auth.Auth) {
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

	ctx, cancel := context.WithCancel(context.Background())

	// Needs to be referenced inside the cancel button callback, to hide the dialog
	var loadingDialog *widget.PopUp

	loadingLabel := widget.NewLabel("Logging in..  Continue in your browser")
	loadingSpinner := widget.NewProgressBarInfinite()
	manualLink := widget.NewHyperlink("If the browser didn't open, click here to login manually", authInstance.LoginURL)
	cancelButton := widget.NewButton("Cancel", func() {
		logrus.Info("Login cancelled by user")
		// This will complete the context, causing WaitForToken to stop waiting for a token.
		cancel()
		loadingDialog.Hide()
	})

	loadingContent := container.NewVBox(loadingLabel, loadingSpinner, manualLink, cancelButton)
	loadingDialog = widget.NewModalPopUp(loadingContent, getScreenManager().window.Canvas())
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
		defer cancel()
		token, err := authInstance.WaitForToken(ctx)
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
		logrus.Info("Login complete")
		getScreenManager().HandleLoginSuccess(token)
	}()
}
