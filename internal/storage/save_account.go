package storage

import (
	"encoding/json"
	"go_diffie_hellman_chat/internal/models"
	"os"
)

// SaveAccount ...
func SaveAccount(account *models.Account, filename string) error {
	jsonData, err := json.Marshal(account)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonData, 0600)
	if err != nil {
		return err
	}

	return nil
}
