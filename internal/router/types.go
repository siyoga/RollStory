package router

import (
	"context"
	"sync"

	uuid "github.com/satori/go.uuid"
	"github.com/siyoga/rollstory/internal/adapter/telegram"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/logger"
)

// Command enums
// naming rules: if command supposed to used via replymarkup keyboard -- must have Plain postfix, if it is default slash command use Slash postfix
type Command string

const (
	StartSlash  Command = "/start"
	CancelSlash Command = "/cancel" // default command which abort current executing command

	WorldSlash Command = "/world"
	WorldPlain Command = "🌍Игровой мир"

	CharacterSlash Command = "/character"
	CharacterPlain Command = "🙎🏻‍♂️Персонаж"

	EditCharacterSlash Command = "/edit_character"
	EditCharacterPlain Command = "🙎🏻‍♂️Изменить описание персонажа"

	EditWorldSlash Command = "/edit_world"
	EditWorldPlain Command = "🌍Изменить описание игрового мира"

	BeginSlash Command = "/begin"
	BeginPlain Command = "▶️Начать игру"

	NewGameSlash Command = "/new_game"
	NewGamePlain Command = "📝Новая игра"
)

// Router and routes
type (
	Router interface {
		Route(triggers ...Command) *route
		DefaultRoute(handler cmdHandler) *route
		Run()
		Stop()
	}

	route struct {
		id      uuid.UUID
		handler cmdHandler

		// linked commands will not cancel each other when called
		linked []Command
	}

	execRoute struct {
		route  *route
		userId int

		respText      string
		respMessageId int
		respChatId    int
	}

	router struct {
		debug bool

		mu           sync.Mutex
		shutdownChan chan struct{}

		triggers     map[Command]uuid.UUID
		routes       map[uuid.UUID]*route
		routesExec   map[int]*execRoute // mapping to already executing routes by user
		defaultRoute *route

		client telegram.Adapter

		ctxHandler  Handler
		gameHandler Handler

		logger logger.Logger
	}
)

type (
	Handler interface {
		FillHandlers(r Router)
	}

	cmdHandler func(ctx context.Context, userId int, msg domain.Request) response

	response struct {
		result domain.Request
		error  *errors.Error
		// pass the user ID of the person whose execution of the command has finished
		release *int
	}
)
