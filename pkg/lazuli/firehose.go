package lazuli

import (
	"bytes"
	"context"
	"io"
	"net/http"

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
		return newError(http.StatusInternalServerError, "fail to connect to websocket", err.Error())
	}
	defer conn.Close()

	for {
		_, message, errMessage := conn.ReadMessage()
		if errMessage != nil {
			return newError(http.StatusInternalServerError, "fail to read message from websocket", errMessage.Error())
		}

		decoder := cbor.NewDecoder(bytes.NewReader(message))

		for {
			var evt dto.RepoCommitEvent
			decodeErr := decoder.Decode(&evt)
			if decodeErr != nil {
				if decodeErr == io.EOF {
					break
				}
				return newError(http.StatusInternalServerError, "fail to decode repo commit event message", decodeErr.Error())
			}

			if handleErr := handler(evt); handleErr != nil {
				return handleErr
			}
		}
	}
}