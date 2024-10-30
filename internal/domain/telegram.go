package domain

// update as needed ref: https://core.telegram.org/bots/api#messageentity
type SpecialType string
type InternalCommand string
type RequestType int

const (
	Mention SpecialType = "mention"
	Hashtag SpecialType = "hashtag"
	Cashtag SpecialType = "cashtag"
	Command SpecialType = "bot_command"
	Url     SpecialType = "url"
	Email   SpecialType = "email"
)

const (
	Cancel InternalCommand = "/cancel"
)

const (
	Message RequestType = iota
	Callback
)

// Markups a.k.a keyboard/buttons
type (
	Markup interface {
		AddRow(row []Button)
	}

	InlineMarkup struct {
		Keyboard map[int][]Button //represent better version of 2d-array map[1] - first row and go on
	}

	ReplyMarkup struct {
		Keyboard     map[int][]Button
		IsPersistent bool
		Resize       bool
		OneTime      bool
	}

	Button struct {
		Text string
		Data string // for inline markups
	}
)

func (im InlineMarkup) AddRow(row []Button) {
	im.Keyboard[len(im.Keyboard)+1] = row
}

func (rm ReplyMarkup) AddRow(row []Button) {
	rm.Keyboard[len(rm.Keyboard)+1] = row
}

type (
	User struct {
		Id       int
		Username string
	}

	Special struct {
		Type   SpecialType
		Offset int
		Length int
	}

	Request struct {
		Id   int
		Type RequestType

		CallbackId *string // pointer cause not every request is a callback.
		MessageId  int
		ChatId     int

		ReplyTo *int // message_id reply to
		From    User
		Data    string

		Command  *string
		Specials []Special
		Markup   Markup
	}
)
