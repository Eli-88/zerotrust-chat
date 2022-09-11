package logger

import (
	"fmt"
	"log"
	"os"
)

// interface compliance
var _ Logger = &consoleLogger{}

type consoleLogger struct {
	level  Level
	logger *log.Logger
}

func (c consoleLogger) internalLog(level Level, args ...any) {
	currentLevel := c.level.Get()
	if currentLevel <= level {
		msg := fmt.Sprintf("[%s] %s", level.ToString(), fmt.Sprintln(args...))
		c.logger.Output(4, msg)
		if level == FATAL {
			panic(msg)
		}
	}
}

func (c *consoleLogger) SetLevel(level Level) {
	c.level.Set(level)
}

func (c consoleLogger) Trace(args ...any) {
	c.internalLog(TRACE, args...)
}

func (c consoleLogger) Debug(args ...any) {
	c.internalLog(DEBUG, args...)
}

func (c consoleLogger) Info(args ...any) {
	c.internalLog(INFO, args...)
}

func (c consoleLogger) Warn(args ...any) {
	c.internalLog(WARN, args...)
}

func (c consoleLogger) Error(args ...any) {
	c.internalLog(ERROR, args...)
}

func (c consoleLogger) Fatal(args ...any) {
	c.internalLog(FATAL, args...)
}

func makeConsoleLogger(level Level) Logger {
	return &consoleLogger{
		level:  level,
		logger: log.New(os.Stdout, "", log.Llongfile),
	}
}
