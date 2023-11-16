package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"go_diffie_hellman_chat/internal/config"
	"go_diffie_hellman_chat/internal/services"
	"go_diffie_hellman_chat/internal/storage"
	"go_diffie_hellman_chat/internal/ui"
)

func main() {
	c := config.DefaultConfig()
	application := app.New()
	window := application.NewWindow("DiffHell")
	window.Resize(fyne.NewSize(400, 500))

	account, err := storage.LoadAccount("./account.json")
	if err != nil {
		//dialog.ShowError(err, window)
		dialog.ShowInformation(
			"Create account",
			"If you already have a private key for Ethereum network,\nyou can create an account with this key,\nthen all your messages will be displayed in your account",
			window,
		)
	}

	ui.ShowAccountScreen(c, window, account, services.CreateAccount)

	window.ShowAndRun()
}
