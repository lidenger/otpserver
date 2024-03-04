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

func InitLog() {
	conf := server.Conf.Log
	atomicLevel, err := zap.ParseAtomicLevel(conf.Level)
	if err != nil {
		panic(fmt.Sprintf("日志Level设置错误:%+v", err))
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConf()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(rotatedConf())), // 打印到控制台和文件
		atomicLevel,
	)
	// 堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	zapLogger = zap.New(core, caller, development)
	Info("日志配置完成")
}

func Info(template string, args ...any) {
	zapLogger.Sugar().Infow(template, args)
}

func Warn(template string, args ...any) {
	zapLogger.Sugar().Warnw(template, args)
}

func Error(template string, args ...any) {
	zapLogger.Sugar().Errorw(template, args)
}

// rotatedConf 日志翻滚切割配置
func rotatedConf() *lumberjack.Logger {
	logConf := server.Conf.Log
	rotatedWriter := &lumberjack.Logger{
		Filename:   logConf.RootPath + logConf.AppFile,
		MaxSize:    logConf.MaxSize,
		MaxBackups: logConf.MaxBackups,
		MaxAge:     logConf.MaxAge,
		Compress:   logConf.Compress,
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
