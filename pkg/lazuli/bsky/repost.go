package bsky

import "time"

type RepostRecord struct {
	LexiconTypeID string        `json:"$type"`
	Subject       RepoStrongRef `json:"subject"`
	CreatedAt     time.Time     `json:"createdAt"`
}
