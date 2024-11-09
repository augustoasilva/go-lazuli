package lazuli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli/bsky"
	"github.com/gorilla/websocket"
)

type Client interface {
	ConsumeFirehose(ctx context.Context, handler HandlerCommitFn) error
	CreateSession(ctx context.Context, identifier, password string) (*bsky.AuthResponse, error)
	CreatePostRecord(ctx context.Context, p bsky.CreateRecordParams) error
	CreateRepostRecord(ctx context.Context, p bsky.CreateRecordParams) error
	CreateLikeRecord(ctx context.Context, p bsky.CreateRecordParams) error
	GetPosts(ctx context.Context, atURIs ...string) (bsky.Posts, error)
	GetPost(ctx context.Context, atURI string) (*bsky.Post, error)
}

type client struct {
	xrpcURL    string
	wsURL      string
	wsDialer   *websocket.Dialer
	session    *bsky.AuthResponse
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

func (c *client) createRecord(ctx context.Context, p bsky.CreateRecordParams) error {
	body := bsky.RequestRecordBody{
		LexiconTypeID: p.Resource,
		Collection:    p.Resource,
		Repo:          c.session.DID,
		Record: bsky.RequestRecord{
			Text: p.Text,
			Subject: bsky.RepoStrongRef{
				URI: p.URI,
				CID: p.CID,
			},
			CreatedAt: time.Now().UTC(),
		},
	}

	jsonBody, _ := json.Marshal(body)

	reqURL := fmt.Sprintf("%s/com.atproto.repo.createRecord", c.xrpcURL)
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonBody))
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

func (c *client) CreatePostRecord(ctx context.Context, p bsky.CreateRecordParams) error {
	p.Resource = "app.bsky.feed.post"
	return c.createRecord(ctx, p)
}

func (c *client) CreateRepostRecord(ctx context.Context, p bsky.CreateRecordParams) error {
	p.Resource = "app.bsky.feed.repost"
	return c.createRecord(ctx, p)
}

func (c *client) CreateLikeRecord(ctx context.Context, p bsky.CreateRecordParams) error {
	p.Resource = "app.bsky.feed.like"
	return c.createRecord(ctx, p)
}

func (c *client) GetPosts(ctx context.Context, atURIs ...string) (bsky.Posts, error) {
	if len(atURIs) == 0 {
		return nil, newError(http.StatusBadRequest, "invalid uris query param", "uris must have at least one value")
	}
	if len(atURIs) > 25 {
		return nil, newError(http.StatusBadRequest, "invalid uris query param", "uris must have at most 25 values")
	}

	query := url.Values{
		"uris": atURIs,
	}
	query.Encode()
	reqURL := fmt.Sprintf("%s/app.bsky.feed.getPosts?%s", c.xrpcURL, query.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, newError(http.StatusInternalServerError, "fail to create get posts request struct", err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.session.AccessJwt))
	req.Header.Set("Content-Type", "application/json")

	resp, doErr := c.httpClient.Do(req)
	if doErr != nil {
		return nil, newError(http.StatusInternalServerError, "fail to do request to get posts", doErr.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, newErrorFromResponse(resp, "get posts request failed")
	}

	var postsResponse bsky.PostResponse
	if decodeErr := json.NewDecoder(resp.Body).Decode(&postsResponse); decodeErr != nil {
		return nil, newError(http.StatusInternalServerError, "fail to decode get posts response", decodeErr.Error())
	}

	return postsResponse.Posts, nil
}

func (c *client) GetPost(ctx context.Context, atURI string) (*bsky.Post, error) {
	posts, err := c.GetPosts(ctx, atURI)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, newError(http.StatusNotFound, "post not found", "post not found")
	}
	return &posts[0], nil
}
