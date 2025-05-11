package models

type User struct {
	ID                int    `json:"id"`
	TelegramUsername  string `json:"telegram_username"`
	TelegramChatID    string `json:"telegram_chat_id"`
	TelegramFirstName string `json:"telegram_first_name"`
	TelegramLastName  string `json:"telegram_last_name"`
}
