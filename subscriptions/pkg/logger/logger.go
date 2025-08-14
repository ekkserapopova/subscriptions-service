package logger

import (
	"log/slog"
	"os"
)

const logFile = "subscriptionService.log"

func SetupLogger() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	return log
}
