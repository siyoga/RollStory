package deps

import (
	"github.com/siyoga/rollstory/internal/adapter"
	"github.com/siyoga/rollstory/internal/adapter/gpt"
	"github.com/siyoga/rollstory/internal/adapter/telegram"
)

func (d *dependencies) GptAdapter() adapter.OpenAIAdapter {
	if d.gptAdapter == nil {
		d.gptAdapter = gpt.NewAdapter(d.cfg.OpenAI, d.log)
	}

	return d.gptAdapter
}

func (d *dependencies) TelegramAdapter() adapter.TelegramAdapter {
	if d.telegramAdapter == nil {
		d.telegramAdapter = telegram.NewAdapter(d.cfg.Bot, d.log)
	}

	return d.telegramAdapter
}
