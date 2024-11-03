package lazuli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli/bsky"
)

func (c *client) CreateSession(ctx context.Context, identifier, password string) (*bsky.AuthResponse, error) {
	request := bsky.SessionRequest{
		Identifier: identifier,
		Password:   password,
	}
	requestBody, _ := json.Marshal(request)

	reqURL := fmt.Sprintf("%s/com.atproto.server.createSession", c.xrpcURL)

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, newError(http.StatusInternalServerError, "fail to create session request struct", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, newError(http.StatusInternalServerError, "error to create session", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, newErrorFromResponse(resp, "create session request failed")
	}

	var didResponse bsky.AuthResponse
	if jsonDecoderErr := json.NewDecoder(resp.Body).Decode(&didResponse); jsonDecoderErr != nil {
		return nil, newError(http.StatusInternalServerError, "error to decode json", jsonDecoderErr.Error())
	}

	c.session = &didResponse

	return &didResponse, nil
}
