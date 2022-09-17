package chat

import (
	"encoding/json"
	"errors"
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/logger"
)

var _ ChatReaderWriter = &chatReaderWriter{}
var _ ChatReaderWriterFactory = chatReaderWriterFactory{}

type chatReaderWriter struct {
	secretKey aes.Key
	conn      Conn
}

type chatReaderWriterFactory struct{}

func NewChatReaderFactory() ChatReaderWriterFactory {
	return &chatReaderWriterFactory{}
}

func NewChatReaderWriter(secretKey aes.Key, conn Conn) ChatReaderWriter {
	return &chatReaderWriter{
		secretKey: secretKey,
		conn:      conn,
	}
}

func (c chatReaderWriterFactory) Create(secretKey aes.Key, conn Conn) ChatReaderWriter {
	return NewChatReaderWriter(secretKey, conn)
}

// TODO: integrate with conn read
func (c chatReaderWriter) Read(msg []byte) ([]ChatMessage, error) {
	logger.Debug("read msg:", string(msg))
	var result []ChatMessage
	for len(msg) > 2 {
		// parse chat message header
		id := msg[0]
		msgLen := int(msg[1])
		logger.Debug("id:", id, " msglen:", msgLen)

		if id != 0x01 {
			break
		}

		msg = msg[2:] // offset slice to start of chat body

		// check len of chat body
		if len(msg) < msgLen {
			break
		}

		// decode chat message
		chatMessage := ChatMessage{}
		if err := json.Unmarshal(msg[:msgLen], &chatMessage); err != nil {
			logger.Debug(err)
			break
		}
		decryptData, err := c.secretKey.Decrypt(chatMessage.Data)
		if err != nil {
			logger.Debug(err)
			break
		}
		chatMessage.Data = string(decryptData)
		result = append(result, chatMessage)

		msg = msg[msgLen:] // offset slice to next chat header
	}
	logger.Debug("result:", result)

	if result == nil {
		return nil, errors.New("invalid chat messages")
	}
	return result, nil
}

func (c *chatReaderWriter) Write(msg []byte) error {
	cipherText, err := c.secretKey.Encrypt(msg)
	if err != nil {
		logger.Debug(err)
		return err
	}
	chatMessage := ChatMessage{
		Data: cipherText,
	}

	chatMsg, err := json.Marshal(chatMessage)
	if err != nil {
		logger.Debug(err)
		return err
	}

	response := []byte{0x01, byte(len(chatMsg))}
	response = append(response, chatMsg...)

	logger.Debug("chat writer send:", string(response))
	_, err = c.conn.Write(response)
	return err
}
