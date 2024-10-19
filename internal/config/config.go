package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/spf13/viper"
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
		StoryDSN  string
		ThreadDSN string
		CertLoc   string
	}

	OpenAI struct {
		Token      string
		Assistants []string
	}

	Timeouts struct {
		RequestTimeout time.Duration
	}
)

func loadTelegramSource(v *viper.Viper) (Bot, error) {
	token, host := v.GetString("telegram.token"), v.GetString("telegram.host")
	if token != "" {
		return Bot{
			Token:     token,
			Debug:     v.GetBool("debug"),
			Host:      host,
			Offset:    0,
			BatchSize: 1,
		}, nil
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

func loadRedisSource(v *viper.Viper) (Redis, error) {
	threadDSN, storyDSN := v.GetString("redis.thread_dsn"), v.GetString("redis.story_dsn")
	if threadDSN != "" && storyDSN != "" {
		return Redis{ThreadDSN: threadDSN, StoryDSN: storyDSN, CertLoc: ""}, nil
	}

	path := v.GetString("apis.redis")
	data, err := os.ReadFile(path)
	if err != nil {
		return Redis{}, err
	}

	var redisCreds struct {
		ThreadDSN string `json:"thread_dsn"`
		StoryDSN  string `json:"story_dsn"`
		CertLoc   string `json:"cert_loc"`
	}

	if err := json.Unmarshal(data, &redisCreds); err != nil {
		return Redis{}, err
	}

	return Redis{
		ThreadDSN: redisCreds.ThreadDSN,
		StoryDSN:  redisCreds.StoryDSN,
		CertLoc:   redisCreds.CertLoc,
	}, nil
}

func loadOpenAiInfo(v *viper.Viper) (OpenAI, error) {
	token, assistants := v.GetString("openai.token"), v.GetStringSlice("openai.assistants")
	if token != "" && len(assistants) > 0 {
		return OpenAI{
			Token:      token,
			Assistants: assistants,
		}, nil
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

	tgSource, err := loadTelegramSource(v)
	if err != nil {
		return nil, err
	}

	openAiCreds, err := loadOpenAiInfo(v)
	if err != nil {
		return nil, err
	}

	redisSource, err := loadRedisSource(v)
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
