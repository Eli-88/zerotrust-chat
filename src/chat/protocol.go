package chat

type KeyExchangeRequest struct {
	PubKey string `json:"public_key"`
}

type KeyExchangeResponse struct {
	SecretKey string `json:"secret_key"`
}

type startConnectionRequest struct {
	Id string `json:"id"`
}

type ChatMessage struct {
	Data string `json:"msg"`
}
