package telegram

import (
	"encoding/json"
	"fmt"

	"github.com/siyoga/rollstory/internal/errors"
)

func (a *adapter) SendMessage(req Request) (err error) {
	defer func() {
		if err != nil {
			a.logger.AdapterError(err, errors.ErrTelegramSendMessage)
		}
	}()

	r, err := json.Marshal(req)
	if err != nil {
		return err
	}

	res, err := a.makeRequest(methodSendMessage, r)

	var message Message
	if err := json.Unmarshal(res, &message); err != nil {
		fmt.Println("unmarshal err:", err)
		return err
	}

	fmt.Println("message sended:", message)

	if err != nil {
		return err
	}

	return nil
}
