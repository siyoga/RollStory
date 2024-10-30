package telegram

import (
	"fmt"
	"strings"

	"github.com/forPelevin/gomoji"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/helpers"
)

func (_ Response) FromDomain(req domain.Request) Response {
	var res Response

	res = Response{
		ChatId: req.ChatId,
		Text:   req.Data,
	}

	switch req.Markup.(type) {
	case domain.InlineMarkup:
		markup := req.Markup.(domain.InlineMarkup)
		res.ReplyMarkup = InlineKeyboardMarkup{
			Buttons: InlineKeyboardMarkup{}.FromDomain(markup.Keyboard),
		}

	case domain.ReplyMarkup:
		markup := req.Markup.(domain.ReplyMarkup)
		res.ReplyMarkup = ReplyKeyboardMarkup{
			Buttons:      ReplyKeyboardMarkup{}.FromDomain(markup.Keyboard),
			IsPersistent: markup.IsPersistent,
			Resize:       markup.Resize,
			OneTime:      markup.OneTime,
		}
	}

	if req.ReplyTo != nil {
		res.Reply = ReplyParamaters{
			MessageId: *req.ReplyTo,
		}
	}

	return res
}

func convertButtons(buttons map[int][]domain.Button) [][]Button {
	keyboard := make([][]Button, len(buttons))
	for i := range keyboard {
		row := make([]Button, len(buttons[i+1]))

		for i, button := range buttons[i+1] {
			if button.Data == "" {
				row[i] = Button{
					Text: button.Text,
				}
			} else {
				row[i] = Button{
					Text:         button.Text,
					CallbackData: button.Data,
				}
			}
		}

		keyboard[i] = row
	}

	return keyboard
}

func (_ InlineKeyboardMarkup) FromDomain(keyboard map[int][]domain.Button) [][]Button {
	return convertButtons(keyboard)
}

func (_ ReplyKeyboardMarkup) FromDomain(keyboard map[int][]domain.Button) [][]Button {
	return convertButtons(keyboard)
}

func (m Message) IsCommand() bool {
	if m.Entities != nil && len(m.Entities) != 0 {
		entity := m.Entities[0]
		return entity.Offset == 0 && entity.Type == "bot_command"
	} else {
		return isCommand(m.Text)
	}
}

// is callback sending a command
func (c Callback) IsCommand() bool {
	return isCommand(c.Data)
}

func (m Message) GetCommand() string {
	if !m.IsCommand() {
		return "text"
	} else {
		return m.Text
	}
}

func isCommand(text string) bool {
	if strings.HasPrefix(text, "/") && len(strings.Split(text, " ")) == 1 {
		return true
	}

	if gomoji.ContainsEmoji(text) {
		emojis := gomoji.CollectAll(text)
		if len(emojis) == 1 {
			emoji := emojis[0]
			if len(strings.Split(text, emoji.Character)) == 2 {
				return true
			}
		}
	}

	return false

}

func (s MessageEntity) ToDomain() domain.Special {
	return domain.Special{
		Type:   domain.SpecialType(s.Type),
		Offset: s.Offset,
		Length: s.Length,
	}
}

func (im InlineKeyboardMarkup) ToDomain() domain.InlineMarkup {
	keyboard := make(map[int][]domain.Button)
	for i, row := range im.Buttons {
		r := make([]domain.Button, len(row))

		for k, button := range row {
			r[k] = domain.Button{
				Text: button.Text,
				Data: button.CallbackData,
			}
		}

		keyboard[i+1] = r
	}

	return domain.InlineMarkup{
		Keyboard: keyboard,
	}
}

func (u Update) ToDomain() domain.Request {
	var req domain.Request

	if u.Message != nil {
		req = domain.Request{
			Id:         u.ID,
			Type:       domain.Message,
			CallbackId: nil,
			MessageId:  u.Message.Id,
			ChatId:     u.Message.Chat.ID,
			From: domain.User{
				Id:       u.Message.From.Id,
				Username: u.Message.From.Username,
			},
			Data:     u.Message.Text,
			Specials: helpers.ConvertSlice(u.Message.Entities, MessageEntity.ToDomain),
			// messages that come in the update can only be with InlineKeyboardMarkup
			Markup: u.Message.ReplyMarkup.ToDomain(),
		}

		fmt.Println("isCommand?:", u.Message.IsCommand())
		if u.Message.IsCommand() {
			cmdName := u.Message.GetCommand()
			req.Command = &cmdName
		}

		return req
	}

	if u.Callback != nil {
		req = domain.Request{
			Id:         u.ID,
			Type:       domain.Callback,
			CallbackId: &u.Callback.Id,
			MessageId:  u.Callback.Message.Id,
			ChatId:     u.Callback.Message.Chat.ID,
			From: domain.User{
				Id:       u.Callback.From.Id,
				Username: u.Callback.From.Username,
			},
			Data:     u.Callback.Data,
			Specials: helpers.ConvertSlice(u.Callback.Message.Entities, MessageEntity.ToDomain),
		}

		if u.Callback.IsCommand() {
			cmdName := u.Callback.Data
			req.Command = &cmdName
		}

		return req
	}

	return req
}
