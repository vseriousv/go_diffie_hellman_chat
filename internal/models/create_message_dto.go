package models

type CreateMessageDTO struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message []byte `json:"message"`
}
