package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli"
)

func main() {
	ctx := context.Background()
	xrpcURL := os.Getenv("XRPC_URL")
	wsURL := os.Getenv("WS_URL")

	client := lazuli.NewClient(xrpcURL, wsURL)

	identifier := os.Getenv("IDENTIFIER")
	password := os.Getenv("PASSWORD")
	sess, err := client.CreateSession(ctx, identifier, password)
	if err != nil {
		slog.Error("error creating session", "error", err)
		panic(err)
	}
	slog.Info("session created", "session", sess)
}
