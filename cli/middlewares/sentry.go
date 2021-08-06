package middlewares

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func SentryInit() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://3d429c316dae4194a3d498fb76de0d49@o517696.ingest.sentry.io/5895770",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	// sentry.CaptureMessage("It works!")
}
