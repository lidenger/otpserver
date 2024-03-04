package log

import (
	"fmt"
	"github.com/lidenger/otpserver/config/server"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var zapLogger *zap.Logger

func InitLog(conf *server.Config) {
	level, err := zap.ParseAtomicLevel(conf.Log.Level)
	if err != nil {
		panic(fmt.Sprintf("日志Level设置错误:%+v", err))
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConf()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(rotatedConf(conf))), // 打印到控制台和文件
		level,
	)
	// 堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	zapLogger = zap.New(core, caller, development)
	Info("日志配置完成,日志文件:%s", conf.Log.RootPath+conf.Log.AppFile)
}

func Info(template string, args ...any) {
	zapLogger.Sugar().Infof(template, args)
}

func Warn(template string, args ...any) {
	zapLogger.Sugar().Infof(template, args)
}

func Error(template string, args ...any) {
	zapLogger.Sugar().Infof(template, args)
}

// rotatedConf 日志翻滚切割配置
func rotatedConf(conf *server.Config) *lumberjack.Logger {
	rotatedWriter := &lumberjack.Logger{
		Filename:   conf.Log.RootPath + conf.Log.AppFile,
		MaxSize:    conf.Log.MaxSize,
		MaxBackups: conf.Log.MaxBackups,
		MaxAge:     conf.Log.MaxAge,
		Compress:   conf.Log.Compress,
	}
	return rotatedWriter
}

// encoderConf 配置默认值
func encoderConf() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.DateTime))
	}
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.MessageKey = "msg"
	encoderConfig.FunctionKey = "func"
	encoderConfig.StacktraceKey = "stacktrace"
	return encoderConfig
}
