package lazuli

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_newError(t *testing.T) {
	type in struct {
		code    int
		message string
		details string
	}

	type out struct {
		lazuliErr error
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "Given a valid input, When newError is called, Then it should return a new error",
			in: in{
				code:    http.StatusInternalServerError,
				message: "test message",
				details: "test details",
			},
			out: out{
				lazuliErr: &Error{
					Code:    http.StatusInternalServerError,
					Message: "test message",
					Details: "test details",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newError(tt.in.code, tt.in.message, tt.in.details)
			assert.Equal(t, tt.out.lazuliErr, got)
		})
	}
}

func TestError_newErrorFromResponse(t *testing.T) {
	type in struct {
		resp    *http.Response
		message string
	}

	type out struct {
		lazuliErr error
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "Given a valid response input without body, When newErrorFromResponse is called, Then it should return a new error",
			in: in{
				resp: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
				message: "test message",
			},
			out: out{
				lazuliErr: &Error{
					Code:    http.StatusInternalServerError,
					Message: "test message",
					Details: "could not get detail",
				},
			},
		},
		{
			name: "Given a valid response input with body, When newErrorFromResponse is called, Then it should return a new error",
			in: in{
				resp: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{"message": "test message"}`)),
				},
				message: "test message",
			},
			out: out{
				lazuliErr: &Error{
					Code:    http.StatusInternalServerError,
					Message: "test message",
					Details: `{"message": "test message"}`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newErrorFromResponse(tt.in.resp, tt.in.message)
			assert.Equal(t, tt.out.lazuliErr, got)
		})
	}
}
