package lazuli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli/dto"
	"github.com/gorilla/websocket"
)

type Client interface {
	ConsumeFirehose(ctx context.Context, handler HandlerCommitFn) error
	CreateSession(ctx context.Context, identifier, password string) (*dto.AuthResponse, error)
	CreateRepostRecord(ctx context.Context, p dto.CreateRecordParams) error
	CreateLikeRecord(ctx context.Context, p dto.CreateRecordParams) error
}

type client struct {
	xrpcURL  string
	wsURL    string
	wsDialer *websocket.Dialer
	session  *dto.AuthResponse
}

func NewClient(xrpcURL, wsURL string) Client {
	dialer := *websocket.DefaultDialer
	return &client{
		xrpcURL:  xrpcURL,
		wsURL:    wsURL,
		wsDialer: &dialer,
	}
}

func (c *client) createRecord(ctx context.Context, p dto.CreateRecordParams) error {
	body := dto.RequestRecordBody{
		LexiconTypeID: p.Resource,
		Collection:    p.Resource,
		Repo:          c.session.DID,
		Record: dto.RequestRecord{
			Subject: dto.RepoStrongRef{
				URI: p.URI,
				CID: p.CID,
			},
			CreatedAt: time.Now().UTC(),
		},
	}

	jsonBody, _ := json.Marshal(body)

	url := fmt.Sprintf("%s/com.atproto.repo.createRecord", c.xrpcURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		slog.Error("fail to create record request struct", "error", err, "resource", p.Resource)
		return nil
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.session.AccessJwt))
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		slog.Error("fail to do request to create record", "error", err, "resource", p.Resource)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("create record request failed", "status", resp, "resource", p.Resource)
		return fmt.Errorf("create record request failed: %d", resp.StatusCode)
	}

	slog.Info("record created", "resource", p.Resource)

	return nil
}

func (c *client) CreateRepostRecord(ctx context.Context, p dto.CreateRecordParams) error {
	p.Resource = "app.bsky.feed.repost"
	return c.createRecord(ctx, p)
}

func (c *client) CreateLikeRecord(ctx context.Context, p dto.CreateRecordParams) error {
	p.Resource = "app.bsky.feed.like"
	return c.createRecord(ctx, p)
}
