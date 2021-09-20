package initialize

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"user_srv/global"
)

func InitLogger() {
	mode := global.Config.RunMod
	fileName := global.Config.Logger.FileName
	maxSize := global.Config.Logger.MaxSize
	maxBackups := global.Config.Logger.MaxBackups
	maxAge := global.Config.Logger.MaxAge
	compress := global.Config.Logger.Compress

	// 打印错误级别的日志
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	// 打印所有级别的日志
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	var allCore []zapcore.Core

	//输出到终端
	consoleDebugging := zapcore.Lock(os.Stdout)

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	if mode == "debug" {
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority))
	} else {
		hook := lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    maxSize, // megabytes
			MaxBackups: maxBackups,
			MaxAge:     maxAge,   //days
			Compress:   compress, // disabled by default
		}

		fileWriter := zapcore.AddSync(&hook)
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, fileWriter, highPriority))
	}

	core := zapcore.NewTee(allCore...)

	Logger := zap.New(core).WithOptions(zap.AddCaller())

	zap.ReplaceGlobals(Logger)
}
