package entities

type Webhook struct {
	Ok  bool     `json:"ok"`
	Res []Result `json:"result"`
}

type Result struct {
	UpdateId int64         `json:"update_id"`
	Msg      Message       `json:"message"`
	Query    CallBackQuery `json:"callback_query"`
}

type Message struct {
	Id   int64       `json:"message_id"`
	From FromUser    `json:"from"`
	Chat ChatDetails `json:"chat"`
	// Date     string      `json:"date"`
	Text string `json:"text"`
	// Entities []Ent       `json:"entities"`
}

type TelegramUsers struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

type FromUser struct {
	TelegramUsers
	LangCode string `json:"languageCode"`
}

type ChatDetails struct {
	TelegramUsers
	Type string `json:"type"`
}

type Ent struct {
	Offset int16  `json:"offset"`
	Length int16  `json:"length"`
	Type   string `json:"type"`
}

type TelegramReplyMsg struct {
	ChatId   int    `json:"chatId"`
	UserName string `json:"username"`
	Text     string `json:"text"`
}

type ReplyResponse struct {
	Ok  bool               `json:"ok"`
	Res ReplyResultDetails `json:"result"`
}

type ReplyResultDetails struct {
	MessageId int64   `json:"message_id"`
	Msg       Message `json:"message"`
}

type ButtonsInMessage struct {
	InlineKeyboard [][]Button `json:"inline_keyboard"`
}

type Button struct {
	Text         string `json:"text"`
	Url          string `json:"url"`
	CallbackData string `json:"callback_data"`
}

type CallBackQuery struct {
	Id          string      `json:"id"`
	From        FromUser    `json:"from"`
	Chat        ChatDetails `json:"chat"`
	Msg         Message     `json:"message"`
	Data        string      `json:"data"`
	InlineMsgId string      `json:"inline_message_id"`
}
