package deps

import (
	"github.com/siyoga/rollstory/internal/adapter/gpt"
	"github.com/siyoga/rollstory/internal/adapter/telegram"
)

func (d *dependencies) GptAdapter() gpt.Adapter {
	if d.gptAdapter == nil {
		d.gptAdapter = gpt.NewAdapter(d.cfg.OpenAI, d.log)
	}

	return d.gptAdapter
}

func (d *dependencies) TelegramAdapter() telegram.Adapter {
	if d.telegramAdapter == nil {
		d.telegramAdapter = telegram.NewAdapter(d.cfg.Bot, d.log)
	}

	return d.telegramAdapter
}
