package lazuli

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli/dto"
	"github.com/fxamacker/cbor/v2"
)

type HandlerCommitFn func(evt dto.CommitEvent) error

// ConsumeFirehose connects to a websocket, reads messages, decodes them as repo commit events, and processes them using a handler function.
//
// TODO: improve firehose consumer to be more flexible
func (c *client) ConsumeFirehose(ctx context.Context, handler HandlerCommitFn) error {
	conn, _, err := c.wsDialer.Dial(c.wsURL, nil)
	if err != nil {
		slog.Error("fail to connect to websocket", "error", err)
		return err
	}
	defer conn.Close()

	slog.Info("websocket connected", "url", c.wsURL)

	for {
		_, message, errMessage := conn.ReadMessage()
		if errMessage != nil {
			slog.Error("fail to read message from websocket", "error", errMessage)
			return errMessage
		}

		decoder := cbor.NewDecoder(bytes.NewReader(message))

		for {
			var evt dto.RepoCommitEvent
			decodeErr := decoder.Decode(&evt)
			if decodeErr != nil {
				if decodeErr == io.EOF {
					break
				}
				slog.Error("fail to decode repo commit event message", "error", decodeErr)
				return decodeErr
			}

			if handleErr := handler(evt); handleErr != nil {
				panic(handleErr)
			}
		}
	}
}
