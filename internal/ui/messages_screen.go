package ui

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"go_diffie_hellman_chat/internal/config"
	"go_diffie_hellman_chat/internal/models"
	"go_diffie_hellman_chat/internal/services"
	"time"
)

// ShowMessageListScreen ...
func ShowMessageListScreen(c *config.Config, window fyne.Window, account *models.Account) {
	welcomeMessage := "Welcome, " + account.Name + "!"

	welcomeLabel := widget.NewLabel(welcomeMessage)

	copyPublicKey := widget.NewButton("Copy PublicKey", func() {
		window.Clipboard().SetContent(account.PublicKey)
	})

	box := container.NewVBox()
	box.Add(welcomeLabel)
	box.Add(copyPublicKey)
	box.Add(widget.NewButton("New message", func() {
		ShowCreateMessageScreen(c, window, account)
	}))

	messageBox := container.NewVBox()

	refreshMessages := func() {
		messages := services.GetMessagesByPublicId(c, account.PublicKey)

		for _, m := range messages {
			var chatName string
			var companion string

			if m.From == account.PublicKey {
				companion = m.To
				chatName = fmt.Sprintf("%d. Me => %s", m.Id, companion)
			} else {
				companion = m.From
				chatName = fmt.Sprintf("%d. %s => Me", m.Id, companion)
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

	bodyBox := container.NewVBox(box, messageBox)

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

	content := container.NewVBox(
		toEntry,
		messageEntry,
		sendButton,
		backButton,
	)

	window.SetContent(content)
}

func ShowMessage(c *config.Config, window fyne.Window, account *models.Account, companion string, msg models.MessageDTO) {

	messageResult, err := services.DecryptMessage(account, companion, msg)
	if err != nil {
		dialog.ShowError(err, window)
	}

	companionLabel := widget.NewLabel(fmt.Sprintf("Chat with: %s", companion))

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
