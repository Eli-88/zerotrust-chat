package chat

import (
	"encoding/binary"
	"encoding/json"
	"testing"
	"zerotrust_chat/chat"
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/logger"
	"zerotrust_chat/test/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReadSingleChatMessage(t *testing.T) {
	logger.SetLogLevel(logger.DEBUG)
	mockCtrl := gomock.NewController(t)
	mockConn := mocks.NewMockHandshakeConn(mockCtrl)

	secretKey, _ := aes.GenerateKey()

	plaintext := "Hello Test"
	cipherText, _ := secretKey.Encrypt([]byte(plaintext))
	chatBody := chat.ChatMessage{
		Data: cipherText,
	}
	msgBody, _ := json.Marshal(chatBody)

	msg := []byte{byte(0x01)}
	msgLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(msgLen, uint32(len(msgBody)))
	msg = append(msg, msgLen...)
	msg = append(msg, msgBody...)

	chatReader := chat.NewChatReaderWriter(secretKey, mockConn)
	chatMessages, err := chatReader.Read(msg)
	assert.NoError(t, err)
	assert.NotNil(t, chatMessages)
	assert.Equal(t, len(chatMessages), 1)
	assert.Equal(t, chatMessages[0].Data, plaintext)
}

func TestReadDoubleChatMessageOnlyFirstValid(t *testing.T) {
	logger.SetLogLevel(logger.DEBUG)
	mockCtrl := gomock.NewController(t)
	mockConn := mocks.NewMockHandshakeConn(mockCtrl)
	secretKey, _ := aes.GenerateKey()

	plainText := "Hello Test"
	cipherText, _ := secretKey.Encrypt([]byte(plainText))
	// valid first message setup
	chatBody := chat.ChatMessage{
		Data: cipherText,
	}
	msgBody, _ := json.Marshal(chatBody)
	msgLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(msgLen, uint32(len(msgBody)))
	msg := []byte{byte(0x01)}
	msg = append(msg, msgLen...)
	msg = append(msg, msgBody...)

	// invalid second message setup
	msg = append(msg, []byte{0x01, 100}...)
	msg = append(msg, []byte("invalid")...)

	logger.Debug("msg:", string(msg))

	chatReader := chat.NewChatReaderWriter(secretKey, mockConn)
	chatMessages, err := chatReader.Read(msg)
	assert.NoError(t, err)
	assert.NotNil(t, chatMessages)

	assert.Equal(t, len(chatMessages), 1)
	assert.Equal(t, chatMessages[0].Data, plainText)
}

func TestReadMultipleChatMessage(t *testing.T) {
	logger.SetLogLevel(logger.DEBUG)

	mockCtrl := gomock.NewController(t)
	mockConn := mocks.NewMockHandshakeConn(mockCtrl)

	secretKey, _ := aes.GenerateKey()

	// valid first message setup
	plainText1 := "Hello Test1"
	cipherText1, _ := secretKey.Encrypt([]byte(plainText1))
	chatBody := chat.ChatMessage{
		Data: cipherText1,
	}
	msgBody, _ := json.Marshal(chatBody)
	msgLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(msgLen, uint32(len(msgBody)))
	msg := []byte{byte(0x01)}
	msg = append(msg, msgLen...)
	msg = append(msg, msgBody...)

	// valid second message setup
	plainText2 := "Hello Test2"
	cipherText2, _ := secretKey.Encrypt([]byte(plainText2))
	chatBody2 := chat.ChatMessage{
		Data: cipherText2,
	}
	msgBody2, _ := json.Marshal(chatBody2)

	msg = append(msg, byte(0x01))
	binary.LittleEndian.PutUint32(msgLen, uint32(len(msgBody2)))
	msg = append(msg, msgLen...)
	msg = append(msg, msgBody2...)

	chatReader := chat.NewChatReaderWriter(secretKey, mockConn)
	chatMessages, err := chatReader.Read(msg)
	assert.NoError(t, err)
	assert.NotNil(t, chatMessages)
	assert.Equal(t, len(chatMessages), 2)
	assert.Equal(t, chatMessages[0].Data, plainText1)
	assert.Equal(t, chatMessages[1].Data, plainText2)

}

func TestReadFirstChatMessageFail(t *testing.T) {
	logger.SetLogLevel(logger.DEBUG)
	mockCtrl := gomock.NewController(t)
	mockConn := mocks.NewMockHandshakeConn(mockCtrl)

	secretKey, _ := aes.GenerateKey()
	// invalid first message setup

	msgBody := []byte("invalid")
	msgLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(msgLen, uint32(len(msgBody)))
	msg := []byte{byte(0x01)}
	msg = append(msg, msgLen...)
	msg = append(msg, msgBody...)

	// valid second message setup
	plainText2 := "Hello Test2"
	cipherText, _ := secretKey.Encrypt([]byte(plainText2))
	chatBody2 := chat.ChatMessage{
		Data: cipherText,
	}
	msgBody2, _ := json.Marshal(chatBody2)
	msg = append(msg, []byte{0x01, byte(len(msgBody2))}...)
	msg = append(msg, msgBody2...)

	chatReader := chat.NewChatReaderWriter(secretKey, mockConn)
	chatMessages, err := chatReader.Read(msg)
	assert.Error(t, err)
	assert.Nil(t, chatMessages)
}
