package lazuli

import (
	"context"
	"net/http"
	"testing"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli/dto"
	"github.com/fxamacker/cbor/v2"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
)

// MockHandlerCommitFn is a mock implementation of HandlerCommitFn.
type MockHandlerCommitFn struct {
	mock.Mock
}

func (m *MockHandlerCommitFn) Handle(evt dto.CommitEvent) error {
	args := m.Called(evt)
	return args.Error(0)
}

func TestClient_ConsumeFirehose(t *testing.T) {
	type in struct {
		ctx                    context.Context
		events                 []dto.CommitEvent
		handler                *MockHandlerCommitFn
		shouldBreakCBORDecoder bool
	}

	type out struct {
		err error
	}

	tests := []struct {
		name  string
		in    in
		setup func(events []dto.CommitEvent, handler *MockHandlerCommitFn)
		out   out
	}{
		{
			name: "Given valid events, When ConsumeFirehose is called, Then it should process the events successfully",
			in: in{
				ctx: context.Background(),
				events: []dto.CommitEvent{
					dto.RepoCommitEvent{},
				},
				handler: &MockHandlerCommitFn{},
			},
			setup: func(events []dto.CommitEvent, handler *MockHandlerCommitFn) {
				for _, event := range events {
					handler.On("Handle", event).Return(nil)
				}
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "Given an invalid event, When ConsumeFirehose is called, Then it should return an error",
			in: in{
				ctx: context.Background(),
				events: []dto.CommitEvent{
					dto.RepoCommitEvent{},
				},
				handler:                &MockHandlerCommitFn{},
				shouldBreakCBORDecoder: true,
			},
			setup: func(events []dto.CommitEvent, handler *MockHandlerCommitFn) {
				for _, event := range events {
					handler.On("Handle", event).Return(nil)
				}
			},
			out: out{
				err: newError(http.StatusInternalServerError, "fail to decode repo commit event message", `cbor: unexpected "break" code`),
			},
		},
		{
			name: "Given a handler error, When ConsumeFirehose is called, Then it should return the handler's error",
			in: in{
				ctx: context.Background(),
				events: []dto.CommitEvent{
					dto.RepoCommitEvent{},
				},
				handler: &MockHandlerCommitFn{},
			},
			setup: func(events []dto.CommitEvent, handler *MockHandlerCommitFn) {
				for _, event := range events {
					handler.On("Handle", event).Return(newError(http.StatusInternalServerError, "handler error", "handler error"))
				}
			},
			out: out{
				err: newError(http.StatusInternalServerError, "handler error", "handler error"),
			},
		},
		{
			name: "Given a websocket connection error, When ConsumeFirehose is called, Then it should return a connection error",
			in: in{
				ctx:     context.Background(),
				events:  nil, // No events to process for this test case
				handler: &MockHandlerCommitFn{},
			},
			out: out{
				err: newError(http.StatusInternalServerError, "fail to connect to websocket", "malformed ws or wss URL"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server

			if len(tt.in.events) > 0 {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					upgrader := websocket.Upgrader{}
					conn, err := upgrader.Upgrade(w, r, nil)
					require.NoError(t, err)
					defer conn.Close()

					for _, event := range tt.in.events {
						message, cborErr := cbor.Marshal(event)
						require.NoError(t, cborErr)
						writeErr := conn.WriteMessage(websocket.BinaryMessage, message)
						require.NoError(t, writeErr)
					}

					if tt.in.shouldBreakCBORDecoder {
						err = conn.WriteMessage(websocket.BinaryMessage, []byte{'\xFF'}) // Invalid CBOR to trigger decode error
						require.NoError(t, err)
					}
				}))
			} else {
				server = httptest.NewUnstartedServer(nil)
			}

			defer server.Close()
			wsURL := "ws"
			if len(server.URL) > 4 {
				wsURL = wsURL + server.URL[4:] + "/ws"
			}
			lazuliClient := &client{
				wsURL:    wsURL,
				wsDialer: websocket.DefaultDialer,
			}

			if tt.setup != nil && tt.in.handler != nil {
				tt.setup(tt.in.events, tt.in.handler)
			}

			err := lazuliClient.ConsumeFirehose(tt.in.ctx, tt.in.handler.Handle)

			if tt.out.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.out.err, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.in.handler != nil {
				tt.in.handler.AssertExpectations(t)
			}
		})
	}
}
