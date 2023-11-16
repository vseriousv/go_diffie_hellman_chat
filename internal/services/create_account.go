package services

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/vseriousv/blockchainkeys"
	"go_diffie_hellman_chat/internal/models"
	"go_diffie_hellman_chat/internal/storage"
	"math/big"
	"strings"
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

func GetAddressFromPublicKey(publicKey string) (string, error) {
	publicKeyBytes, err := hex.DecodeString(strings.Split(publicKey, "0x04")[1])
	if err != nil {
		return "", err
	}

	xBytes := publicKeyBytes[:len(publicKeyBytes)/2]
	yBytes := publicKeyBytes[len(publicKeyBytes)/2:]

	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)

	publicKeyECDSA := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	address := crypto.PubkeyToAddress(publicKeyECDSA).Hex()

	return address, nil
}
