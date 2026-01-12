package logger

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

var Log *Logger

func InitLogger(filePath string, fileName string, mode string) error {
	// Ensure the log directory exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, 0777)
		if err != nil {
			return err
		}
	}

	var logger zerolog.Logger

	if mode == "release" {
		logFile, err := os.OpenFile(filepath.Join(filePath, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

		if err != nil {
			return err
		}

		logger = zerolog.New(logFile).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	// Set global level to trace for maximum verbosity (adjust as needed)
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	Log = &Logger{logger}
	return nil
}
