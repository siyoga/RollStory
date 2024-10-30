package telegram

import (
	"encoding/json"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (a *adapter) SendMessage(req Response) (msg Message, err error) {
	defer func() {
		if err != nil {
			a.logger.AdapterError(err, errors.ErrTelegramSendMessage)
		}
	}()

	r, err := json.Marshal(req)
	if err != nil {
		return Message{}, err
	}

	res, err := a.makeRequest(methodSendMessage, r)
	if err != nil {
		return Message{}, err
	}

	var resp struct {
		Ok     bool    `json:"ok"`
		Result Message `json:"result"`
	}
	if err := json.Unmarshal(res, &resp); err != nil {
		return Message{}, err
	}

	return resp.Result, nil
}

func (a *adapter) DeleteMessage(chatId int, messageId int) (err error) {
	defer func() {
		if err != nil {
			a.logger.AdapterError(err, errors.ErrTelegramSendMessage)
		}
	}()

	var req struct {
		ChatId    int `json:"chat_id"`
		MessageId int `json:"message_id"`
	}

	req.ChatId = chatId
	req.MessageId = messageId

	r, err := json.Marshal(req)
	if err != nil {
		return err
	}

	res, err := a.makeRequest(methodDeleteMessage, r)
	if err != nil {
		return err
	}

	var message Message
	if err := json.Unmarshal(res, &message); err != nil {
		return err
	}

	return nil
}

func (a *adapter) EditMessage(chatId int, messageId int, text string, markups *domain.InlineMarkup) (err error) {
	defer func() {
		if err != nil {
			a.logger.AdapterError(err, errors.ErrTelegramSendMessage)
		}
	}()

	var req struct {
		ChatId      int                  `json:"chat_id"`
		MessageId   int                  `json:"message_id"`
		Text        string               `json:"text"`
		ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}

	req.ChatId = chatId
	req.MessageId = messageId
	req.Text = text

	if markups != nil {
		req.ReplyMarkup = InlineKeyboardMarkup{
			Buttons: InlineKeyboardMarkup{}.FromDomain(markups.Keyboard),
		}
	}

	r, err := json.Marshal(req)
	if err != nil {
		return err
	}

	res, err := a.makeRequest(methodEditMessage, r)
	if err != nil {
		return err
	}

	var message Message
	if err := json.Unmarshal(res, &message); err != nil {
		return err
	}

	return nil
}
