package deps

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/siyoga/rollstory/internal/adapter"
	"github.com/siyoga/rollstory/internal/app/db"
	"github.com/siyoga/rollstory/internal/config"
	"github.com/siyoga/rollstory/internal/handler"
	"github.com/siyoga/rollstory/internal/handler/router"
	"github.com/siyoga/rollstory/internal/logger"
	"github.com/siyoga/rollstory/internal/repository"
	"github.com/siyoga/rollstory/internal/service"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Dependencies interface {
		Cfg() *config.Config

		Router() router.Router
		WaitForInterrupt()
		Close()
	}

	dependencies struct {
		shutdownChannel chan os.Signal
		closeCallbacks  []func()

		log         logger.Logger
		cfg         *config.Config
		redisClient *db.RedisClient

		threadRepository repository.ThreadRepository

		contextService service.ContextService
		gameService    service.GameService

		router         router.Router
		contextHandler handler.Handler
		gameHandler    handler.Handler

		gptAdapter      adapter.OpenAIAdapter
		telegramAdapter adapter.TelegramAdapter
	}
)

func NewDependencies(cfgPath string) (Dependencies, error) {
	cfg, err := config.NewConfig(cfgPath)
	if err != nil && err.Error() == "Config File \"config\" Not Found in \"[]\"" {
		cfg, err = config.NewConfig("./configs/local")
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.LevelKey = "lvl"
	encoderCfg.TimeKey = "t"
	l := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.Lock(os.Stdout), zap.NewAtomicLevel()),
		zap.AddCaller(),
	)

	return &dependencies{
		cfg:             cfg,
		log:             logger.NewLogger(l, "rollstory_bot", "github.com/siyoga/rollstory"),
		shutdownChannel: make(chan os.Signal),
	}, nil
}

func (d *dependencies) Cfg() *config.Config {
	return d.cfg
}

func (d *dependencies) Close() {
	for i := len(d.closeCallbacks) - 1; i >= 0; i-- {
		d.closeCallbacks[i]()
	}
}

func (d *dependencies) WaitForInterrupt() {
	signal.Notify(d.shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
	d.log.Zap().Info("Wait for receive interrupt signal...")
	<-d.shutdownChannel
	d.log.Zap().Info("Receive interrupt signal")
}
