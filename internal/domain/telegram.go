package domain

type (
	UpdatesChan <-chan Update

	UpdatesResponse struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	Update struct {
		ID      int      `json:"update_id"`
		Message *Message `json:"message"`
	}

	Message struct {
		Text     string          `json:"text"`
		From     From            `json:"from"`
		Chat     Chat            `json:"chat"`
		Entities []MessageEntity `json:"entities,omitempty"`
	}

	MessageEntity struct {
		Type   string `json:"type"`
		Offset int    `json:"offset"`
		Length int    `json:"length"`
	}

	MessageResult struct {
		Message string
		ChatId  int64
	}

	From struct {
		Id       int64  `json:"int64"`
		Username string `json:"username"`
	}

	Chat struct {
		ID int64 `json:"id"`
	}
)

func (m Message) IsCommand() bool {
	if m.Entities == nil || len(m.Entities) == 0 {
		return false
	}

	entity := m.Entities[0]
	return entity.Offset == 0 && entity.Type == "bot_command"
}

func (m Message) GetCommand() string {
	if !m.IsCommand() {
		return "text"
	}

	entity := m.Entities[0]
	return m.Text[1:entity.Length]
}
