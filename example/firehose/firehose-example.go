package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli"
	lazulidto "github.com/augustoasilva/go-lazuli/pkg/lazuli/bsky"
)

func main() {

	ctx := context.Background()
	xrpcURL := os.Getenv("XRPC_URL")
	wsURL := os.Getenv("WS_URL")

	client := lazuli.NewClient(xrpcURL, wsURL)

	handler := func(evt lazulidto.CommitEvent) error {
		slog.Info("reading 1 firehose event", "type", evt.Type())
		return nil
	}
	err := client.ConsumeFirehose(ctx, handler)
	if err != nil {
		slog.Error("error consuming firehose", "error", err)
		panic(err)
	}
}
