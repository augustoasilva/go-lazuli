package dto

import "time"

type Reply struct {
	Parent RepoStrongRef `json:"parent"`
	Root   RepoStrongRef `json:"root"`
}

type PostRecord struct {
	LexiconTypeID string `json:"$type"`
	Text          string `json:"text"`
	Reply         *Reply `json:"reply,omitempty"`
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
