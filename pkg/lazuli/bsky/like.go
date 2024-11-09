package bsky

import "time"

type LikeRecord struct {
	LexiconTypeID string        `json:"$type"`
	Subject       RepoStrongRef `json:"subject"`
	CreatedAt     time.Time     `json:"createdAt"`
}
