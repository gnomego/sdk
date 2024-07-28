package log

import (
	"fmt"
	"log/slog"
)

func Error(err error, message string, args ...interface{}) {
	msg := message
	if len(args) > 0 {
		msg = fmt.Sprintf(message, args...)
	}

	slog.Error(msg, slog.Any("error", err))
}
