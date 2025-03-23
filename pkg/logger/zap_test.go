package logger_test

import (
	"testing"

	"github.com/AugustineAurelius/fuufu/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWithDebug(t *testing.T) {
	log := logger.New()
	require.Equal(t, log.Level(), zap.InfoLevel)

	log2 := logger.New(logger.WithDebug())
	require.Equal(t, log2.Level(), zap.DebugLevel)
}

func TestWithJson(t *testing.T) {
	log := logger.New()
	log.Info("test log")

	log2 := logger.New(logger.WithJSON())
	log2.Info("test log json")
}
