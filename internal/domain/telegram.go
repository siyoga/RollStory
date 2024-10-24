package domain

import "github.com/siyoga/rollstory/internal/adapter/telegram"

type (
	Message struct {
		UpdateId int
		Text     string
		ChatId   int64
		Command  *command
	}

	command struct {
		name string
	}
)

func (m Message) ToRequest(rows ...[]telegram.Button) telegram.Request {
	keyboard := make([][]telegram.Button, len(rows))

	for i, row := range rows {
		keyboard[i] = row
	}

	return telegram.Request{
		ChatId: m.ChatId,
		Text:   m.Text,
		ReplyMarkup: telegram.ReplyKeyboardMarkup{
			Buttons:      keyboard,
			IsPersistent: true,
		},
	}
}

func (_ Message) FromUpdate(u telegram.Update) Message {
	return Message{
		UpdateId: u.ID,
		Text:     u.Message.Text,
		ChatId:   u.Message.Chat.ID,
		Command:  getCommand(u.Message),
	}
}

func getCommand(m *telegram.Message) *command {
	if m.IsCommand() {
		return &command{
			name: m.GetCommand(),
		}
	}

	return nil
}
