package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Level string

const (
	LevelProd Level = "prod"
	LevelWarn Level = "warn"
	LevelDev  Level = "dev"
)

type Config struct {
	Level          Level
	FileLog        FileLog
	DevelopmerMode bool
}

type FileLog struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Localtime  bool
	Compress   bool
}

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	DPanicw(msg string, args ...interface{})
	New(name string) Logger
}

type zapConfig struct {
	zap        *zap.Config
	lumberjack *lumberjack.Logger
}

func getZapConfig(config Config) zapConfig {
	zaps := zapConfig{
		zap: &zap.Config{
			DisableCaller:     false,
			DisableStacktrace: false,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stdout"},
			EncoderConfig: zapcore.EncoderConfig{
				NameKey:      "name",
				LevelKey:     "level",
				TimeKey:      "time",
				MessageKey:   "msg",
				CallerKey:    "caller",
				EncodeTime:   zapcore.ISO8601TimeEncoder,
				EncodeLevel:  zapcore.LowercaseLevelEncoder,
				EncodeCaller: zapcore.ShortCallerEncoder,
			},
		},
		lumberjack: nil,
	}

	switch config.Level {
	case LevelProd:
		zaps.zap.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		zaps.zap.DisableStacktrace = true
		zaps.zap.DisableStacktrace = true
		zaps.zap.Encoding = "json"
		if config.FileLog.Filename == "" {
			zaps.lumberjack = &lumberjack.Logger{
				Filename:   "./holavonat.log",
				MaxSize:    50, // megabytes
				MaxBackups: 30,
				MaxAge:     28, // days
				LocalTime:  false,
				Compress:   false,
			}
		}
	case LevelWarn:
		zaps.zap.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		zaps.zap.Encoding = "console"
	case LevelDev:
		zaps.zap.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zaps.zap.Encoding = "console"
	default:
		zaps.zap.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		zaps.zap.Encoding = "console"
	}

	if config.FileLog.Filename != "" && zaps.lumberjack == nil {
		zaps.lumberjack = &lumberjack.Logger{
			Filename:   config.FileLog.Filename,
			MaxSize:    50, // megabytes (default)
			MaxBackups: 30, // default
			MaxAge:     28, // days (default)
			LocalTime:  false,
			Compress:   false,
		}
	} else if config.FileLog.Filename != "" && zaps.lumberjack != nil {
		if config.FileLog.MaxSize != 0 {
			zaps.lumberjack.MaxSize = config.FileLog.MaxSize
		}
		if config.FileLog.MaxBackups != 0 {
			zaps.lumberjack.MaxBackups = config.FileLog.MaxBackups
		}
		if config.FileLog.MaxAge != 0 {
			zaps.lumberjack.MaxAge = config.FileLog.MaxAge
		}
		if config.FileLog.Localtime {
			zaps.lumberjack.LocalTime = config.FileLog.Localtime
		}
		if config.FileLog.Compress {
			zaps.lumberjack.Compress = config.FileLog.Compress
		}
	}

	return zaps
}

func checkFilePathValid(fp string) error {
	if _, err := os.Stat(fp); err == nil {
		return err
	}

	var d []byte
	if err := os.WriteFile(fp, d, 0600); err == nil {
		os.Remove(fp)
		return err
	}

	return nil
}

type logger struct {
	zapLogger *zap.Logger
}

func newLogger(config Config) *logger {
	if err := checkFilePathValid(config.FileLog.Filename); err != nil {
		panic(err)
	}

	zapParsedConf := getZapConfig(config)

	if zapParsedConf.lumberjack == nil {
		zapLogger, err := zapParsedConf.zap.Build(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
		logger := logger{zapLogger}
		return &logger
	}

	zapLogger, err := zapParsedConf.zap.Build(zap.WrapCore(zapCore(zapParsedConf.lumberjack)), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	logger := logger{zapLogger}
	return &logger

}

func zapCore(lj *lumberjack.Logger) func(c zapcore.Core) zapcore.Core {
	return func(c zapcore.Core) zapcore.Core {
		w := zapcore.AddSync(lj)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w,
			zap.InfoLevel,
		)
		cores := zapcore.NewTee(c, core)

		return cores
	}
}

func (l *logger) New(name string) Logger {
	newLogger := logger{l.zapLogger.Named(name)}
	return &newLogger
}

func (l *logger) Debug(args ...interface{}) {
	l.zapLogger.Sugar().Debug(args...)
}
func (l *logger) Debugf(template string, args ...interface{}) {
	l.zapLogger.Sugar().Debugf(template, args...)
}
func (l *logger) Debugw(msg string, args ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.zapLogger.Sugar().Info(args...)
}
func (l *logger) Infof(template string, args ...interface{}) {
	l.zapLogger.Sugar().Infof(template, args...)
}
func (l *logger) Infow(msg string, args ...interface{}) {
	l.zapLogger.Sugar().Infow(msg, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.zapLogger.Sugar().Warn(args...)
}
func (l *logger) Warnf(template string, args ...interface{}) {
	l.zapLogger.Sugar().Warnf(template, args...)
}
func (l *logger) Warnw(msg string, args ...interface{}) {
	l.zapLogger.Sugar().Warnw(msg, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.zapLogger.Sugar().Error(args...)
}
func (l *logger) Errorf(template string, args ...interface{}) {
	l.zapLogger.Sugar().Errorf(template, args...)
}
func (l *logger) Errorw(msg string, args ...interface{}) {
	l.zapLogger.Sugar().Errorw(msg, args...)
}

func (l *logger) DPanic(args ...interface{}) {
	l.zapLogger.Sugar().DPanic(args...)
}
func (l *logger) DPanicf(template string, args ...interface{}) {
	l.zapLogger.Sugar().DPanicf(template, args...)
}
func (l *logger) DPanicw(msg string, args ...interface{}) {
	l.zapLogger.Sugar().DPanicw(msg, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.zapLogger.Sugar().Fatal(args...)
}
func (l *logger) Fatalf(template string, args ...interface{}) {
	l.zapLogger.Sugar().Fatalf(template, args...)
}
func (l *logger) Fatalw(msg string, args ...interface{}) {
	l.zapLogger.Sugar().Fatalw(msg, args...)
}
