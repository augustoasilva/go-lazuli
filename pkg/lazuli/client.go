package lazuli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli/dto"
	"github.com/gorilla/websocket"
)

type Client interface {
	ConsumeFirehose(ctx context.Context, handler HandlerCommitFn) error
	CreateSession(ctx context.Context, identifier, password string) (*dto.AuthResponse, error)
	CreatePostRecord(ctx context.Context, p dto.CreateRecordParams) error
	CreateRepostRecord(ctx context.Context, p dto.CreateRecordParams) error
	CreateLikeRecord(ctx context.Context, p dto.CreateRecordParams) error
}

type client struct {
	xrpcURL    string
	wsURL      string
	wsDialer   *websocket.Dialer
	session    *dto.AuthResponse
	httpClient *http.Client
}

func NewClient(xrpcURL, wsURL string) Client {
	dialer := *websocket.DefaultDialer
	// TODO: improve to use a more appropriate http client config
	return &client{
		xrpcURL:    xrpcURL,
		wsURL:      wsURL,
		wsDialer:   &dialer,
		httpClient: http.DefaultClient,
	}
}

func (c *client) createRecord(ctx context.Context, p dto.CreateRecordParams) error {
	body := dto.RequestRecordBody{
		LexiconTypeID: p.Resource,
		Collection:    p.Resource,
		Repo:          c.session.DID,
		Record: dto.RequestRecord{
			Text: p.Text,
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
		return newError(http.StatusInternalServerError, "fail to create record request struct", err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.session.AccessJwt))
	req.Header.Set("Content-Type", "application/json")

	resp, doErr := c.httpClient.Do(req)
	if doErr != nil {
		return newError(http.StatusInternalServerError, "fail to do request to create record", doErr.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return newErrorFromResponse(resp, "create record request failed")
	}

	return nil
}

func (c *client) CreatePostRecord(ctx context.Context, p dto.CreateRecordParams) error {
	p.Resource = "app.bsky.feed.post"
	return c.createRecord(ctx, p)
}

func (c *client) CreateRepostRecord(ctx context.Context, p dto.CreateRecordParams) error {
	p.Resource = "app.bsky.feed.repost"
	return c.createRecord(ctx, p)
}

func (c *client) CreateLikeRecord(ctx context.Context, p dto.CreateRecordParams) error {
	p.Resource = "app.bsky.feed.like"
	return c.createRecord(ctx, p)
}
