package logger

var root *logger

func init() {
	root = newLogger(Config{})
}
func Init(config Config) {
	root = newLogger(config)
}

func Debug(args ...interface{}) {
	root.zapLogger.Sugar().Debug(args...)
}
func Debugf(template string, args ...interface{}) {
	root.zapLogger.Sugar().Debugf(template, args...)
}
func Debugw(msg string, args ...interface{}) {
	root.zapLogger.Sugar().Debugw(msg, args...)
}

func Info(args ...interface{}) {
	root.zapLogger.Sugar().Info(args...)
}
func Infof(template string, args ...interface{}) {
	root.zapLogger.Sugar().Infof(template, args...)
}
func Infow(msg string, args ...interface{}) {
	root.zapLogger.Sugar().Infow(msg, args...)
}

func Warn(args ...interface{}) {
	root.zapLogger.Sugar().Warn(args...)
}
func Warnf(template string, args ...interface{}) {
	root.zapLogger.Sugar().Warnf(template, args...)
}
func Warnw(msg string, args ...interface{}) {
	root.zapLogger.Sugar().Warnw(msg, args...)
}

func Error(args ...interface{}) {
	root.zapLogger.Sugar().Error(args...)
}
func Errorf(template string, args ...interface{}) {
	root.zapLogger.Sugar().Errorf(template, args...)
}
func Errorw(msg string, args ...interface{}) {
	root.zapLogger.Sugar().Errorw(msg, args...)
}

func DPanic(args ...interface{}) {
	root.zapLogger.Sugar().DPanic(args...)
}
func DPanicf(template string, args ...interface{}) {
	root.zapLogger.Sugar().DPanicf(template, args...)
}
func DPanicw(msg string, args ...interface{}) {
	root.zapLogger.Sugar().DPanicw(msg, args...)
}

func Fatal(args ...interface{}) {
	root.zapLogger.Sugar().Fatal(args...)
}
func Fatalf(template string, args ...interface{}) {
	root.zapLogger.Sugar().Fatalf(template, args...)
}
func Fatalw(msg string, args ...interface{}) {
	root.zapLogger.Sugar().Fatalw(msg, args...)
}

func New(name string) Logger {
	return root.New(name)
}
