package cmd

import (
	"github.com/AugustineAurelius/fuufu/internal/config"
	"github.com/AugustineAurelius/fuufu/pkg/logger"
	"go.uber.org/zap"
)

func getLogger(manager *config.Manager) *zap.Logger {
	logOpts := make([]logger.LoggerOpt, 0, 8)
	logCfg := manager.LoadLogging()
	if logCfg.Debug {
		logOpts = append(logOpts, logger.WithDebug())
	}
	if logCfg.Json {
		logOpts = append(logOpts, logger.WithJson())
	}

	return logger.New(logOpts...)
}
