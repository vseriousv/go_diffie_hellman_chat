package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAddressFromPublicKey(t *testing.T) {
	publicKey := "0x040978abde5e642ff14ca7ebc3d5099bc4761e22c05e9d911e2246e5e244001c68f0e3b0c4aa414135bda0f4950d7449a7e63768d14060f31ee8df75544b44c299"
	address := "0x7488fAe195561A1187acafeEDc971Fa56D4eEBb5"

	actual, err := GetAddressFromPublicKey(publicKey)
	assert.Nil(t, err, "error: ", err)
	assert.Equal(t, address, actual)
}
