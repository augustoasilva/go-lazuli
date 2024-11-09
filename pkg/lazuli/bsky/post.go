package bsky

import "time"

// PostAuthor
//
// TODO: add remaining fields associated,viewer and labels
type PostAuthor struct {
	DID         string    `json:"did"`
	Handle      string    `json:"handle"`
	DisplayName string    `json:"displayName"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"createdAt"`
}

// PostViewer
//
//	Metadata about the requesting account's relationship with the subject content. Only has meaningful content for authed requests.
type PostViewer struct {
	Repost            string `json:"repost"` // at-uri
	Like              string `json:"like"`   // at-uri
	ThreadMuted       bool   `json:"threadMuted"`
	ReplyDisabled     bool   `json:"replyDisabled"`
	EmbeddingDisabled bool   `json:"embeddingDisabled"`
	Pinned            bool   `json:"pinned"`
}

// PostRecord
//
// Represents the post record data.
type PostRecord struct {
	LexiconTypeID string           `json:"$type"`
	URI           string           `json:"uri"`
	CreatedAt     time.Time        `json:"createdAt"`
	Embed         any              `json:"embed"` // TODO: embed can be many types of objects, for now it will be any, need improvement
	Facets        []map[string]any `json:"facets,omitempty"`
	Langs         []string         `json:"langs"`
	Reply         *Reply           `json:"reply,omitempty"`
	Text          string           `json:"text"`
}

// Post
//
// Represents the struct with data when requesting the get posts endpoint
type Post struct {
	LexiconTypeID string     `json:"$type"`
	URI           string     `json:"uri"` // at-uri
	CID           any        `json:"cid"` // TODO: for now handle CID like the the firehose, but need improvement
	Author        PostAuthor `json:"author"`
	Record        PostRecord `json:"record"`
	Embed         any        `json:"embed"` // TODO: embed can be many types of objects, for now it will be any, need improvement
	ReplyCount    int        `json:"replyCount,omitempty"`
	RepostCount   int        `json:"repostCount,omitempty"`
	LikeCount     int        `json:"likeCount,omitempty"`
	QuoteCount    int        `json:"quoteCount,omitempty"`
	IndexedAt     time.Time  `json:"indexedAt"`
	Viewer        PostViewer `json:"viewer"`
}

type PostRecords []Post

type PostResponse struct {
	Posts PostRecords `json:"posts"`
}
