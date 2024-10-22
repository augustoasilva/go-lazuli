package lazuli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli/dto"
)

func (c *client) CreateSession(ctx context.Context, identifier, password string) (*dto.AuthResponse, error) {
	request := dto.SessionRequest{
		Identifier: identifier,
		Password:   password,
	}
	requestBody, _ := json.Marshal(request)

	url := fmt.Sprintf("%s/com.atproto.server.createSession", c.xrpcURL)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("create session request failed", "status_code", resp.StatusCode)
		return nil, fmt.Errorf("create session request failed: %d", resp.StatusCode)
	}

	var didResponse dto.AuthResponse
	if jsonDecoderErr := json.NewDecoder(resp.Body).Decode(&didResponse); jsonDecoderErr != nil {
		slog.Error(
			"error to decod json response",
			"error", jsonDecoderErr,
		)
		return nil, jsonDecoderErr
	}

	c.session = &didResponse

	return &didResponse, nil
}
