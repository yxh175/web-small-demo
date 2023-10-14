package logger

import (
	"gin-mall/elk_demo/util/rwriter"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(rw *rwriter.RedisWriter) *zap.Logger {

	// 设置日志级别
	lowPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapcore.DebugLevel
	})

	// 使用json格式日志
	jsonEnc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	stdCore := zapcore.NewCore(jsonEnc, zapcore.Lock(os.Stdout), lowPriority)

	// rw实现io.Writer的接口
	syncer := zapcore.AddSync(rw)
	redisCore := zapcore.NewCore(jsonEnc, syncer, lowPriority)

	// 集成多个内核
	core := zapcore.NewTee(stdCore, redisCore)

	// Logger 输出console且标识调用代码行
	return zap.New(core).WithOptions(zap.AddCaller())
}
