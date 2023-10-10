package telegram

type TgClient interface {
	Updates(offset, limit int) ([]Update, error)
	SendMessage(chatID int, text string) error
	SendPhotoWithInlineKeyboard(chatID int, text string, photoURL string, buttons []InlineKeyboardButton) error
	SendMessageWithInlineKeyboard(chatID int, text string, buttons []InlineKeyboardButton) error
	AnswerCallbackQuery(queryID string, text string) error

	EditMessageReplyMarkup(chatID int, messageID int) error
	DeleteMessage(chatID int, messageID int) (bool, error)
}

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message"`
	CallbackQuery *CallbackQuery   `json:"callback_query"`
}

type IncomingMessage struct {
	ID   int    `json:"message_id"`
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type DeleteResponse struct {
	Result bool `json:"result"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
