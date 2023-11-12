package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vseriousv/diffiehellman/pkg/encryption"
	"github.com/vseriousv/diffiehellman/pkg/transport_key"
	"go_diffie_hellman_chat/internal/config"
	"go_diffie_hellman_chat/internal/models"
	"io"
	"log"
	"net/http"
)

func SendMessage(c *config.Config, msg models.CreateMessageDTO, account *models.Account) (bool, error) {
	url := fmt.Sprintf("%s/v1/messages", c.ApiUrl)

	transportKey, err := transport_key.GetTransportKey(msg.To, account.PrivateKey)
	if err != nil {
		return false, err
	}

	encryptionMessage, err := encryption.Encrypt(msg.Message, []byte(transportKey))
	if err != nil {
		return false, err
	}

	msg.Message = encryptionMessage

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return false, err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	log.Println("Response Status:", response.Status)
	log.Println("Response Body:", string(body))

	return true, nil
}

func DecryptMessage(account *models.Account, companion string, msg models.MessageDTO) (*string, error) {
	transportKey, err := transport_key.GetTransportKey(companion, account.PrivateKey)
	if err != nil {
		return nil, err
	}

	messageResult, err := encryption.Decrypt(msg.Message, []byte(transportKey))
	if err != nil {
		return nil, err
	}

	return &messageResult, nil
}

func GetMessagesByPublicId(c *config.Config, pubKey string) []models.MessageDTO {
	url := fmt.Sprintf("%s/v1/messages/%s", c.ApiUrl, pubKey)

	response, err := http.Get(url)
	if err != nil {
		log.Printf(err.Error())
		return nil
	}

	if err != nil {
		log.Printf("Error happened in sending request. Err: %s", err.Error())
		return nil
	}

	defer response.Body.Close()

	var messages []models.MessageDTO

	err = json.NewDecoder(response.Body).Decode(&messages)
	if err != nil {
		log.Printf("Error happened in sending request. Err: %s", err.Error())
		return nil
	}

	return messages
}
