package chat

import (
	"encoding/json"
	"testing"
	"zerotrust_chat/chat"
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/crypto/rsa"
	"zerotrust_chat/test/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServerHandshake(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockHandshakeConn := mocks.NewMockHandshakeConn(mockCtrl)
	mockKeyFactory := mocks.NewMockKeyFactory(mockCtrl)

	priKey, _ := rsa.GeneratePrivateKey()
	secretKey, _ := aes.GenerateKey()

	keyRequest := chat.KeyExchangeRequest{
		PubKey: priKey.GetPublicKey().ToString(),
	}
	req, _ := json.Marshal(keyRequest)

	encryptedSecretKey, _ := priKey.GetPublicKey().Encrypt([]byte(secretKey.ToString()))
	keyResponse := chat.KeyExchangeResponse{
		SecretKey: encryptedSecretKey,
	}
	res, _ := json.Marshal(keyResponse)

	mockKeyFactory.EXPECT().GenerateRsaPrivateKey().Return(priKey, nil)
	mockHandshakeConn.EXPECT().Write(req).Return(len(req), nil)
	mockHandshakeConn.EXPECT().Read(gomock.Any()).Return(res, nil)
	mockKeyFactory.EXPECT().ConstructAesSecretKey(secretKey.ToString()).Return(secretKey, nil)

	key, err := chat.NewServerHandshake(mockHandshakeConn, mockKeyFactory).Handshake()
	assert.NoError(t, err)
	assert.Equal(t, key.ToString(), secretKey.ToString())
}
