package logger

import (
	"gin-mall/elk_demo/util/rwriter"
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {

	rwriter := rwriter.NewRedisWriter()
	logger := NewLogger(rwriter)

	logger.Info("test logger info", zap.String("hello", "string"))
}
