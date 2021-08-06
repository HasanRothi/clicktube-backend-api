package middlewares

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

func SentryInit() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	// sentry.CaptureMessage("It works!")
}
