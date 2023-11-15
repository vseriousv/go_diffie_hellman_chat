package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"go_diffie_hellman_chat/internal/config"
	"go_diffie_hellman_chat/internal/models"
)

// NewAccountScreen ...
func NewAccountScreen(c *config.Config, window fyne.Window, createFunc func(name string, privateKey *string) (*models.Account, error)) {
	nameEntry := widget.NewEntry()
	privateKeyEntry := widget.NewEntry()
	privateKeyEntry.SetPlaceHolder("Optional")

	createButton := widget.NewButton("Create Account", func() {
		var privateKeyPtr *string
		if privateKeyEntry.Text != "" {
			privateKeyPtr = &privateKeyEntry.Text
		}
		account, err := createFunc(nameEntry.Text, privateKeyPtr)
		if err != nil {
			dialog.ShowError(err, window)
		}
		ShowMessageListScreen(c, window, account)
	})

	form := container.NewVBox(
		widget.NewLabel("Enter account name:"),
		nameEntry,
		widget.NewLabel("Enter private key (if you have one):"),
		privateKeyEntry,
		createButton,
	)

	window.SetContent(form)
}

// ShowAccountScreen ...
func ShowAccountScreen(
	c *config.Config,
	window fyne.Window,
	account *models.Account,
	createFunc func(name string, privateKey *string) (*models.Account, error),
) {
	if account == nil {
		NewAccountScreen(c, window, createFunc)
	} else {
		ShowMessageListScreen(c, window, account)
	}
}
