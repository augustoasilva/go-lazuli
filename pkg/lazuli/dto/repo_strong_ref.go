package dto

type RepoStrongRef struct {
	LexiconTypeID string `json:"$type,omitempty"`
	CID           string `json:"cid"`
	URI           string `json:"uri"`
}
