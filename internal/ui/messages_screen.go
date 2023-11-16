package ui

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"go_diffie_hellman_chat/internal/config"
	"go_diffie_hellman_chat/internal/models"
	"go_diffie_hellman_chat/internal/services"
	"go_diffie_hellman_chat/internal/utils"
	"time"
)

// ShowMessageListScreen ...
func ShowMessageListScreen(c *config.Config, window fyne.Window, account *models.Account) {
	welcomeMessage := "Welcome, " + account.Name + "!"

	welcomeLabel := widget.NewLabel(welcomeMessage)

	copyPublicKeyButton := widget.NewButton("Copy PublicKey", func() {
		window.Clipboard().SetContent(account.PublicKey)
	})

	newMessageButton := widget.NewButton("New message", func() {
		ShowCreateMessageScreen(c, window, account)
	})

	box := container.NewVBox()
	box.Add(welcomeLabel)
	box.Add(container.NewPadded(container.New(layout.NewGridLayout(2), copyPublicKeyButton, newMessageButton)))

	messageBox := container.NewVBox()

	refreshMessages := func() {
		messages := services.GetMessagesByPublicId(c, account.PublicKey)

		for _, m := range messages {
			var chatName string
			var companion string

			messageData := m.CreatedAt.Format("2006-01-02")
			messageTime := m.CreatedAt.Format("15:04")

			if m.From == account.PublicKey {
				companion = m.To
				address, err := services.GetAddressFromPublicKey(m.To)
				if err != nil {
					dialog.ShowError(err, window)
					return
				}
				chatName = fmt.Sprintf("%s %s \t Me => %s", messageData, messageTime, utils.AddressShort(address))
			} else {
				companion = m.From
				address, err := services.GetAddressFromPublicKey(m.From)
				if err != nil {
					dialog.ShowError(err, window)
					return
				}
				chatName = fmt.Sprintf("%s %s \t %s => Me", messageData, messageTime, utils.AddressShort(address))
			}

			currentMessage := m

			messageBox.Add(widget.NewButton(chatName, func() {
				ShowMessage(c, window, account, companion, currentMessage)
			}))
		}
	}

	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			messageBox.RemoveAll()

			refreshMessages()

			window.Content().Refresh()
		}
	}()

	refreshMessages()

	bodyBox := container.NewVBox(
		box,
		container.NewPadded(messageBox),
	)

	window.SetContent(bodyBox)
}

func ShowCreateMessageScreen(c *config.Config, window fyne.Window, account *models.Account) {
	toEntry := widget.NewEntry()

	messageEntry := widget.NewMultiLineEntry()

	sendButton := widget.NewButton("Send", func() {
		msg := models.CreateMessageDTO{
			From:    account.PublicKey,
			To:      toEntry.Text,
			Message: []byte(messageEntry.Text),
		}

		if toEntry.Text == "" || messageEntry.Text == "" {
			dialog.ShowError(errors.New("Fields must be filled in"), window)
			return
		}

		_, err := services.SendMessage(c, msg, account)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		ShowMessageListScreen(c, window, account)
	})

	backButton := widget.NewButton("Back", func() {
		ShowMessageListScreen(c, window, account)
	})

	buttonBox := container.New(layout.NewGridLayout(2), backButton, sendButton)

	content := container.NewVBox(
		widget.NewLabel("To publicKey:"),
		toEntry,
		widget.NewLabel("Message:"),
		messageEntry,
		buttonBox,
	)

	window.SetContent(content)
}

func ShowMessage(c *config.Config, window fyne.Window, account *models.Account, companion string, msg models.MessageDTO) {

	messageResult, err := services.DecryptMessage(account, companion, msg)
	if err != nil {
		dialog.ShowError(err, window)
	}

	address, err := services.GetAddressFromPublicKey(companion)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	companionLabel := widget.NewLabel(fmt.Sprintf("Chat with: %s", address))

	messageTitleLabel := widget.NewLabel("Message:")

	messageLabel := widget.NewLabel(*messageResult)

	backButton := widget.NewButton("Back", func() {
		ShowMessageListScreen(c, window, account)
	})

	content := container.NewVBox(
		companionLabel,
		messageTitleLabel,
		messageLabel,
		backButton,
	)

	window.SetContent(content)
}
