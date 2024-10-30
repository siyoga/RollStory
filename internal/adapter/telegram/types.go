package telegram

import "github.com/siyoga/rollstory/internal/domain"

// Markups
type (
	ReplyMarkup interface {
		FromDomain(keyboard map[int][]domain.Button) [][]Button
	}

	// https://core.telegram.org/bots/api#replykeyboardmarkup
	ReplyKeyboardMarkup struct {
		Buttons      [][]Button `json:"keyboard"`
		IsPersistent bool       `json:"is_persistent"`
		Resize       bool       `json:"resize_keyboard"`
		OneTime      bool       `json:"one_time_keyboard"`
	}

	// https://core.telegram.org/bots/api#inlinekeyboardmarkup
	InlineKeyboardMarkup struct {
		Buttons [][]Button `json:"inline_keyboard"`
	}

	Button struct {
		Text         string `json:"text"`
		CallbackData string `json:"callback_data,omitempty"`
	}
)

// Update
type (
	UpdatesChan <-chan Update

	UpdatesResponse struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	Update struct {
		ID       int       `json:"update_id"`
		Message  *Message  `json:"message,omitempty"`
		Callback *Callback `json:"callback_query,omitempty"`
	}

	UpdateParam struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	}
)

// Update payloads (Message, callbacks e.g.)
type (
	Message struct {
		Id          int                  `json:"message_id"`
		Text        string               `json:"text"`
		From        User                 `json:"from"`
		Chat        Chat                 `json:"chat"`
		Entities    []MessageEntity      `json:"entities,omitempty"`
		ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}

	Callback struct {
		Id      string  `json:"id"`
		From    User    `json:"from"`
		Data    string  `json:"data,omitempty"`
		Message Message `json:"message,omitempty"`
	}

	// basically typings for sendMessage request which appears in BotAPI
	Response struct {
		ChatId      int             `json:"chat_id"`
		Text        string          `json:"text"`
		ReplyMarkup ReplyMarkup     `json:"reply_markup,omitempty"`
		Reply       ReplyParamaters `json:"reply_parameters,omitempty"`
	}

	ReplyParamaters struct {
		MessageId int `json:"message_id"`
	}

	MessageEntity struct {
		Type   string `json:"type"`
		Offset int    `json:"offset"`
		Length int    `json:"length"`
	}

	MessageResult struct {
		Message string
		ChatId  int
	}
)

// Other stuff (from, chat e.g.)
type (
	User struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	}

	Chat struct {
		ID int `json:"id"`
	}
)
