package telegram

import (
	"net/url"
	"strconv"
)

func (a *adapter) SendMessage(chatID int64, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.FormatInt(chatID, 10))
	q.Add("text", text)

	_, err := a.makeRequest(methodSendMessage, q)
	if err != nil {
		return err
	}

	return nil
}
