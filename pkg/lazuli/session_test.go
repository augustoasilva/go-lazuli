package lazuli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli/dto"
	"github.com/stretchr/testify/assert"
)

func TestClient_CreateSession(t *testing.T) {
	type in struct {
		ctx        context.Context
		identifier string
		password   string
	}

	type out struct {
		authResponse *dto.AuthResponse
		err          error
	}

	tests := []struct {
		name    string
		in      in
		out     out
		handler http.HandlerFunc
	}{
		{
			name: "Given valid input and successful response, When CreateSession is called, Then it should return auth response",
			in: in{
				ctx:        context.Background(),
				identifier: "test-user",
				password:   "test-password",
			},
			out: out{
				authResponse: &dto.AuthResponse{AccessJwt: "valid-token"},
				err:          nil,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(dto.AuthResponse{AccessJwt: "valid-token"})
			},
		},
		{
			name: "Given invalid input, When CreateSession is called, Then it should return an error",
			in: in{
				ctx:        context.Background(),
				identifier: "test-user",
				password:   "wrong-password",
			},
			out: out{
				authResponse: nil,
				err: &Error{
					Code:    http.StatusUnauthorized,
					Message: "create session request failed",
					Details: `{"message":"Unauthorized"}` + "\n",
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
			},
		},
		{
			name: "Given server error, When CreateSession is called, Then it should return an internal server error",
			in: in{
				ctx:        context.Background(),
				identifier: "test-user",
				password:   "test-password",
			},
			out: out{
				authResponse: nil,
				err: &Error{
					Code:    http.StatusInternalServerError,
					Message: "error to decode json",
					Details: "unexpected EOF",
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{`)) // invalid JSON
			},
		},
		{
			name: "Given no response from server, When CreateSession is called, Then it should return an internal server error",
			in: in{
				ctx:        context.Background(),
				identifier: "test-user",
				password:   "test-password",
			},
			out: out{
				authResponse: nil,
				err: &Error{
					Code:    http.StatusInternalServerError,
					Message: "error to create session",
					Details: `parse ":invalid-url/com.atproto.server.createSession": missing protocol scheme`,
				},
			},
			handler: nil, // simulate no response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			lazuliClient := &client{
				xrpcURL: server.URL,
			}

			if tt.handler == nil {
				lazuliClient.xrpcURL = ":invalid-url"
			}

			result, err := lazuliClient.CreateSession(tt.in.ctx, tt.in.identifier, tt.in.password)

			if tt.out.err != nil {
				assert.Nil(t, result)
				var actualErr *Error
				assert.ErrorAs(t, err, &actualErr)
				assert.Equal(t, tt.out.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.out.authResponse, result)
			}
		})
	}
}
