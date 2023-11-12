package services

import (
	"github.com/vseriousv/blockchainkeys"
	"go_diffie_hellman_chat/internal/models"
	"go_diffie_hellman_chat/internal/storage"
)

func CreateAccount(name string, privateKey *string) (*models.Account, error) {
	var pk, pb, addr string

	bc, err := blockchainkeys.NewBlockchain("Ethereum")
	if err != nil {
		return nil, err
	}

	if privateKey != nil {
		pk = *privateKey
		pb, addr, err = bc.GetPublicKeyAndAddressByPrivateKey(pk)
		if err != nil {
			return nil, err
		}
	} else {
		pk, pb, addr, err = bc.GenerateKeyPair()
		if err != nil {
			return nil, err
		}
	}

	account := models.Account{
		Name:       name,
		PrivateKey: pk,
		PublicKey:  pb,
		Address:    addr,
	}

	err = storage.SaveAccount(&account, "./account.json")
	if err != nil {
		return nil, err
	}

	return &account, nil
}
