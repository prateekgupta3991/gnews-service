package entities

//User struct
type User struct {
	ID   string `json:"id"`
	Name string `json:"username"`
}

//User Details struct
type UserDetails struct {
	ID         int64  `json:"id"`
	Name       string `json:"username"`
	TelegramId string `json:"telegramId"`
	ChatId     int64  `json:"chatId"`
}
