package chat

import (
	"encoding/json"
	"testing"
	"zerotrust_chat/chat"
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/crypto/rsa"
	test "zerotrust_chat/test/helper"
	"zerotrust_chat/test/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewClientHandshake(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockHandshakeConn := mocks.NewMockHandshakeConn(mockCtrl)
	mockKeyFactory := mocks.NewMockKeyFactory(mockCtrl)

	privateKey, _ := rsa.GeneratePrivateKey()
	secretKey, _ := aes.GenerateKey()

	keyExchangeRequest := test.KeyExchangeRequest{
		PubKey: privateKey.GetPublicKey().ToString(),
	}
	exchangeRequest, _ := json.Marshal(keyExchangeRequest)

	mockHandshakeConn.EXPECT().Write(gomock.Any()).Times(2)
	mockHandshakeConn.EXPECT().Read(gomock.Any()).Return(exchangeRequest, nil)

	mockKeyFactory.EXPECT().ConstructRsaPublicKey(keyExchangeRequest.PubKey).Return(privateKey.GetPublicKey(), nil)
	mockKeyFactory.EXPECT().GenerateAesSecretKey().Return(secretKey, nil)

	key, err := chat.NewClientHandshake(
		"hello",
		mockHandshakeConn,
		mockKeyFactory).Handshake()

	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.True(t, key.ToString() == secretKey.ToString())
}
