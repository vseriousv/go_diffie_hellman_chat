package storage

import (
	"encoding/json"
	"go_diffie_hellman_chat/internal/models"
	"os"
)

// LoadAccount ...
func LoadAccount(filename string) (*models.Account, error) {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var account models.Account
	err = json.Unmarshal(jsonData, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
