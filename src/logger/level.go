package logger

import "sync/atomic"

type Level int32

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

func (l Level) ToString() string {
	switch l {
	case TRACE:
		return "trace"
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	default:
		return "fatal"
	}
}

func (l Level) Get() Level {
	return Level(atomic.LoadInt32((*int32)(&l)))
}

func (l *Level) Set(level Level) {
	atomic.StoreInt32((*int32)(l), (int32)(level))
}
