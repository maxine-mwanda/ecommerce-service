// internal/logging/logging.go
package logging

import (
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rs/zerolog"
)

type RotateLogger struct {
	*zerolog.Logger
	rotateLogs *rotatelogs.RotateLogs
}

func NewRotateLogger(logDir, filename string) (*RotateLogger, error) {
	// Ensure log directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logPath := filepath.Join(logDir, filename)
	rl, err := rotatelogs.New(
		logPath+".%Y%m%d",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	logger := zerolog.New(rl).With().Timestamp().Logger()

	return &RotateLogger{
		Logger:     &logger,
		rotateLogs: rl,
	}, nil
}

func (l *RotateLogger) Close() error {
	return l.rotateLogs.Close()
}
