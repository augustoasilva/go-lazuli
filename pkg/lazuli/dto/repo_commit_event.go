package dto

type RepoCommitEvent struct {
	Repo   string          `cbor:"repo"`
	Rev    string          `cbor:"rev"`
	Seq    int64           `cbor:"seq"`
	Since  string          `cbor:"since"`
	Time   string          `cbor:"time"`
	TooBig bool            `cbor:"tooBig"`
	Prev   any             `cbor:"prev"`
	Rebase bool            `cbor:"rebase"`
	Blocks []byte          `cbor:"blocks"`
	Ops    []RepoOperation `cbor:"ops"`
	Commit any             `json:"commit"` // Repo commit object CID
}

func (e RepoCommitEvent) Type() CommitEventType {
	return CommitEventTypeRepoCommit
}

type RepoOperation struct {
	Action string `cbor:"action"`
	Path   string `cbor:"path"`
	Reply  *Reply `cbor:"reply,omitempty"`
	Text   []byte `cbor:"text"`
	CID    any    `cbor:"cid"`
}
