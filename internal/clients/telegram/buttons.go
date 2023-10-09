package telegram

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type CallbackQuery struct {
	ID      string          `json:"id"`
	Data    string          `json:"data"`
	Message IncomingMessage `json:"message"`
	From    From            `json:"from"`
}
