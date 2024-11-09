package bsky

import "time"

type Record struct {
	CID any    `json:"cid"`
	URI string `json:"uri"`
}

type RequestRecord struct {
	Subject   RepoStrongRef `json:"subject"`
	Text      string        `json:"text,omitempty"`
	CreatedAt time.Time     `json:"createdAt"`
}

type RequestRecordBody struct {
	LexiconTypeID string        `json:"$type"`
	Collection    string        `json:"collection"`
	Repo          string        `json:"repo"`
	Record        RequestRecord `json:"record"`
}

type RequestLikesFromPost struct {
	URI string `json:"uri"`
}

type CreateRecordParams struct {
	Resource string
	Text     string
	URI      string
	CID      string
}
