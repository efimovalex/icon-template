package main

import (
	"fmt"
	"os"

	"github.com/iconimpact/replaceme/config"
	server "github.com/iconimpact/replaceme/services"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// BuildName data
	BuildName = "replaceme"
	// BuildDate data
	BuildDate string
	// BuildBranch data
	BuildBranch string
	// BuildNumber data
	BuildNumber string
)

// main - main entry point that loads configuration and starts the services
func main() {
	// bootstrap logger
	var err error

	var cfg *config.Config

	// load config
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		cfg, err = config.Load(configFile)
	} else {
		cfg, err = config.Load(configFile)
	}
	if err != nil {
		fmt.Printf("config load error %s", err.Error())
		os.Exit(1)
	}

	logConfig := zap.NewProductionConfig()
	logConfig.Level, err = zap.ParseAtomicLevel(cfg.Logging.Level)
	logConfig.Development = cfg.Logging.Development
	logConfig.EncoderConfig.MessageKey = "message"
	logConfig.EncoderConfig.TimeKey = "timestamp"
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := logConfig.Build()
	if err != nil {
		fmt.Printf("failed to parse log level: %s", err.Error())
		os.Exit(1)
	}
	log := logger.Sugar()
	defer log.Sync()

	log.Infow("config loaded", "config", cfg)

	// start services
	server, err := server.New(cfg, log)
	if err != nil {
		fmt.Printf("failed init service: %s", err.Error())
		os.Exit(1)
	}
	// Optional steps: run migrations, seed data, etc. here.

	// start server
	os.Exit(server.Start())
}