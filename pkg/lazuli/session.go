package lazuli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
		return nil, newError(http.StatusInternalServerError, "error to create session", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, newErrorFromResponse(resp, "create session request failed")
	}

	var didResponse dto.AuthResponse
	if jsonDecoderErr := json.NewDecoder(resp.Body).Decode(&didResponse); jsonDecoderErr != nil {
		return nil, newError(http.StatusInternalServerError, "error to decode json", jsonDecoderErr.Error())
	}

	c.session = &didResponse

	return &didResponse, nil
}
