package chat

type keyExchangeRequest struct {
	PubKey string `json:"public_key"`
}

type keyExchangeResponse struct {
	SecretKey string `json:"secret_key"`
}

type startConnectionRequest struct {
	Id string `json:"id"`
}

type ChatMessage struct {
	Data string `json:"msg"`
}
