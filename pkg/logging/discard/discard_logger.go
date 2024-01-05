package discard

import (
	"speed-typing-text-service/pkg/logging"
)

var l *DiscardsLogger

type DiscardsLogger struct {
}

func GetLogger() logging.Logger {
	return l
}

func init() {
	l = &DiscardsLogger{}
}

func (l *DiscardsLogger) Print(_ ...interface{}) {
}

func (l *DiscardsLogger) Printf(_ string, _ ...interface{}) {
}

func (l *DiscardsLogger) Trace(_ ...interface{}) {
}

func (l *DiscardsLogger) Tracef(_ string, _ ...interface{}) {
}

func (l *DiscardsLogger) Debug(_ ...interface{}) {
}

func (l *DiscardsLogger) Debugf(_ string, _ ...interface{}) {
}

func (l *DiscardsLogger) Info(_ ...interface{}) {
}

func (l *DiscardsLogger) Infof(_ string, _ ...interface{}) {
}

func (l *DiscardsLogger) Warn(_ ...interface{}) {
}

func (l *DiscardsLogger) Warnf(_ string, _ ...interface{}) {
}

func (l *DiscardsLogger) Error(_ ...interface{}) {
}

func (l *DiscardsLogger) Errorf(_ string, _ ...interface{}) {
}

func (l *DiscardsLogger) Fatal(_ ...interface{}) {
}

func (l *DiscardsLogger) Fatalf(_ string, _ ...interface{}) {
}

func (l *DiscardsLogger) Panic(_ ...interface{}) {
}

func (l *DiscardsLogger) Panicf(_ string, _ ...interface{}) {
}
