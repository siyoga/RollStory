package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type (
	Mode string
)

const (
	Local Mode = "local"
	Dev   Mode = "dev"
	Prod  Mode = "prod"
)

type (
	Config struct {
		Bot      Bot
		OpenAI   OpenAI
		Timeouts Timeouts
		Redis    Redis
	}

	Bot struct {
		Token string
		Debug bool
		Host  string

		Offset    int
		BatchSize int
	}

	Redis struct {
		DSN     string
		CertLoc string
	}

	OpenAI struct {
		Token      string
		Assistants []string
	}

	Timeouts struct {
		RequestTimeout time.Duration
	}
)

func loadTelegramSource(v *viper.Viper, mode Mode) (Bot, error) {
	token, host := v.GetString("telegram.token"), v.GetString("telegram.host")
	if mode == Local {
		if token != "" && host != "" {
			return Bot{
				Token:     token,
				Debug:     v.GetBool("debug"),
				Host:      host,
				Offset:    0,
				BatchSize: 1,
			}, nil
		} else {
			return Bot{}, fmt.Errorf("telegram.token and telegram.host must be specified with local mode")
		}
	}

	path := v.GetString("apis.telegram")
	data, err := os.ReadFile(path)
	if err != nil {
		return Bot{}, err
	}

	var creds struct {
		Token string `json:"token"`
		Host  string `json:"host"`
	}

	if err := json.Unmarshal(data, &creds); err != nil {
		return Bot{}, err
	}

	return Bot{
		Token:     creds.Token,
		Host:      creds.Host,
		Debug:     false,
		Offset:    0,
		BatchSize: 1,
	}, nil
}

func loadRedisSource(v *viper.Viper, mode Mode) (Redis, error) {
	DSN := v.GetString("redis.dsn")
	if mode == Local {
		if DSN != "" {
			return Redis{DSN: DSN, CertLoc: ""}, nil
		} else {
			return Redis{}, fmt.Errorf("redis.thread_dsn and redis.story_dsn must be specified with local mode")
		}
	}

	path := v.GetString("apis.redis")
	data, err := os.ReadFile(path)
	if err != nil {
		return Redis{}, err
	}

	var redisCreds struct {
		DSN     string `json:"dsn"`
		CertLoc string `json:"cert_loc"`
	}

	if err := json.Unmarshal(data, &redisCreds); err != nil {
		return Redis{}, err
	}

	return Redis{
		DSN:     redisCreds.DSN,
		CertLoc: redisCreds.CertLoc,
	}, nil
}

func loadOpenAiInfo(v *viper.Viper, mode Mode) (OpenAI, error) {
	token, assistants := v.GetString("openai.token"), v.GetStringSlice("openai.assistants")
	if mode == Local {
		if token != "" && len(assistants) > 0 {
			return OpenAI{
				Token:      token,
				Assistants: assistants,
			}, nil
		} else {
			return OpenAI{}, fmt.Errorf("openai.token and openai.assistants must be specified with local mode")
		}
	}

	path := v.GetString("apis.openai")
	data, err := os.ReadFile(path)
	if err != nil {
		return OpenAI{}, err
	}

	var openAICreds struct {
		Token      string   `json:"token"`
		Assistants []string `json:"assistants"`
	}

	if err := json.Unmarshal(data, &openAICreds); err != nil {
		return OpenAI{}, err
	}

	return OpenAI{
		Token:      token,
		Assistants: assistants,
	}, nil
}

func NewConfig(cfgPath string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(cfgPath)
	v.SetConfigName("config")
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	mode := Mode(v.GetString("mode"))

	tgSource, err := loadTelegramSource(v, mode)
	if err != nil {
		return nil, err
	}

	openAiCreds, err := loadOpenAiInfo(v, mode)
	if err != nil {
		return nil, err
	}

	redisSource, err := loadRedisSource(v, mode)
	if err != nil {
		return nil, err
	}

	return &Config{
		Bot:    tgSource,
		OpenAI: openAiCreds,
		Timeouts: Timeouts{
			RequestTimeout: v.GetDuration("timeouts.request"),
		},
		Redis: redisSource,
	}, nil
}
