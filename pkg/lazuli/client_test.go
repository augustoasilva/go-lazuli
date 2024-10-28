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

func TestClient_CreatePostRecord(t *testing.T) {
	type in struct {
		ctx    context.Context
		params dto.CreateRecordParams
	}

	type out struct {
		err error
	}

	tests := []struct {
		name    string
		in      in
		out     out
		handler http.HandlerFunc
	}{
		{
			name: "Given a CreatePostRecord function call, When there is valid params and successful response, Then it should create a post record",
			in: in{
				ctx: context.Background(),
				params: dto.CreateRecordParams{
					Text:     "test text",
					URI:      "test-uri",
					CID:      "test-cid",
					Resource: "app.bsky.feed.post",
				},
			},
			out: out{
				err: nil,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			name: "Given a CreatePostRecord function call, When there is request creation failure, Then it should return an error",
			in: in{
				ctx: context.Background(),
				params: dto.CreateRecordParams{
					Text:     "test text",
					URI:      "test-uri",
					CID:      "test-cid",
					Resource: "app.bsky.feed.post",
				},
			},
			out: out{
				err: newError(http.StatusInternalServerError, "create record request failed", `{"message":"request failed"}`+"\n"),
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"message": "request failed"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			lazuliClient := &client{
				xrpcURL:    server.URL,
				session:    &dto.AuthResponse{AccessJwt: "test-token", DID: "test-did"},
				httpClient: server.Client(),
			}

			err := lazuliClient.CreatePostRecord(tt.in.ctx, tt.in.params)

			if tt.out.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.out.err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_GetPosts(t *testing.T) {
	type in struct {
		ctx    context.Context
		atURIs []string
	}

	type out struct {
		posts dto.Posts
		err   error
	}

	tests := []struct {
		name    string
		in      in
		out     out
		handler http.HandlerFunc
	}{
		{
			name: "Given a GetPosts function call, When there is valid atURIs and successful response, Then it should return posts",
			in: in{
				ctx:    context.Background(),
				atURIs: []string{"test-uri-1", "test-uri-2"},
			},
			out: out{
				posts: dto.Posts{
					{URI: "test-uri-1"},
					{URI: "test-uri-2"},
				},
				err: nil,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				posts := dto.Posts{
					{URI: "test-uri-1"},
					{URI: "test-uri-2"},
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(posts)
			},
		},
		{
			name: "Given a GetPosts function call, When there is a request failure, Then it should return an error",
			in: in{
				ctx:    context.Background(),
				atURIs: []string{"test-uri-1", "test-uri-2"},
			},
			out: out{
				posts: nil,
				err:   newError(http.StatusInternalServerError, "get posts request failed", `{"message":"request failed"}`+"\n"),
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"message": "request failed"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			lazuliClient := &client{
				xrpcURL:    server.URL,
				session:    &dto.AuthResponse{AccessJwt: "test-token"},
				httpClient: server.Client(),
			}

			posts, err := lazuliClient.GetPosts(tt.in.ctx, tt.in.atURIs...)

			if tt.out.err != nil {
				assert.Nil(t, posts)
				assert.Error(t, err)
				assert.Equal(t, tt.out.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.out.posts, posts)
			}
		})
	}
}

func TestClient_GetPost(t *testing.T) {
	type in struct {
		ctx   context.Context
		atURI string
	}

	type out struct {
		post *dto.Post
		err  error
	}

	tests := []struct {
		name    string
		in      in
		out     out
		handler http.HandlerFunc
	}{
		{
			name: "Given a GetPost function call, When there is a valid atURI and successful response, Then it should return the post",
			in: in{
				ctx:   context.Background(),
				atURI: "test-uri",
			},
			out: out{
				post: &dto.Post{
					URI: "test-uri",
				},
				err: nil,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				post := dto.Post{URI: "test-uri"}
				posts := dto.Posts{post}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(posts)
			},
		},
		{
			name: "Given a GetPost function call, When there is an invalid atURI, Then it should return an error",
			in: in{
				ctx:   context.Background(),
				atURI: "invalid-uri",
			},
			out: out{
				post: nil,
				err:  newError(http.StatusInternalServerError, "get posts request failed", `{"message":"request failed"}`+"\n"),
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"message": "request failed"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			lazuliClient := &client{
				xrpcURL:    server.URL,
				session:    &dto.AuthResponse{AccessJwt: "test-token"},
				httpClient: server.Client(),
			}

			post, err := lazuliClient.GetPost(tt.in.ctx, tt.in.atURI)

			if tt.out.err != nil {
				assert.Nil(t, post)
				assert.Error(t, err)
				assert.Equal(t, tt.out.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.out.post, post)
			}
		})
	}
}

func TestClient_CreateRepostRecord(t *testing.T) {
	type in struct {
		ctx    context.Context
		params dto.CreateRecordParams
	}
	type out struct {
		err error
	}
	tests := []struct {
		name    string
		in      in
		out     out
		handler http.HandlerFunc
	}{
		{
			name: "Given a CreateRepostRecord function call, when valid params are provided, then it should create a repost record and return no error",
			in: in{
				ctx: context.Background(),
				params: dto.CreateRecordParams{
					Text:     "test text",
					URI:      "test-uri",
					CID:      "test-cid",
					Resource: "app.bsky.feed.repost",
				},
			},
			out: out{
				err: nil,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			name: "Given a CreateRepostRecord function call, when request creation fails, then it should return an error",
			in: in{
				ctx: context.Background(),
				params: dto.CreateRecordParams{
					Text:     "test text",
					URI:      "test-uri",
					CID:      "test-cid",
					Resource: "app.bsky.feed.repost",
				},
			},
			out: out{
				err: newError(http.StatusInternalServerError, "create record request failed", `{"message":"request failed"}`+"\n"),
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"message": "request failed"})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			lazuliClient := &client{
				xrpcURL:    server.URL,
				session:    &dto.AuthResponse{AccessJwt: "test-token", DID: "test-did"},
				httpClient: server.Client(),
			}
			err := lazuliClient.CreateRepostRecord(tt.in.ctx, tt.in.params)
			if tt.out.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.out.err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_CreateLikeRecord(t *testing.T) {
	type in struct {
		ctx    context.Context
		params dto.CreateRecordParams
	}
	type out struct {
		err error
	}
	tests := []struct {
		name    string
		in      in
		out     out
		handler http.HandlerFunc
	}{
		{
			name: "Given a CreateLikeRecord function call, when valid params are provided, it should create a like record and return no error",
			in: in{
				ctx: context.Background(),
				params: dto.CreateRecordParams{
					Text:     "test text",
					URI:      "test-uri",
					CID:      "test-cid",
					Resource: "app.bsky.feed.like",
				},
			},
			out: out{
				err: nil,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			name: "Given a CreateLikeRecord function call, when request creation fails, then it should return an error",
			in: in{
				ctx: context.Background(),
				params: dto.CreateRecordParams{
					Text:     "test text",
					URI:      "test-uri",
					CID:      "test-cid",
					Resource: "app.bsky.feed.like",
				},
			},
			out: out{
				err: newError(http.StatusInternalServerError, "create record request failed", `{"message":"request failed"}`+"\n"),
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"message": "request failed"})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			lazuliClient := &client{
				xrpcURL:    server.URL,
				session:    &dto.AuthResponse{AccessJwt: "test-token", DID: "test-did"},
				httpClient: server.Client(),
			}
			err := lazuliClient.CreateLikeRecord(tt.in.ctx, tt.in.params)
			if tt.out.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.out.err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
