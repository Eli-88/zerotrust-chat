package logger

var printer = makeConsoleLogger(INFO)

// not thread safe
func SetLogger(logger Logger) {
	printer = logger
}

func Info(args ...any) {
	printer.Info(args...)
}

func Debug(args ...any) {
	printer.Debug(args...)
}

func Trace(args ...any) {
	printer.Trace(args...)
}

func Warn(args ...any) {
	printer.Warn(args...)
}

func Error(args ...any) {
	printer.Error(args...)
}

func Fatal(args ...any) {
	printer.Fatal(args...)
}

func SetLogLevel(level Level) {
	printer.SetLevel(level)
}
