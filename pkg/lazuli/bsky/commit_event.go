package bsky

type CommitEvent interface {
	Type() CommitEventType
	GetOps() []RepoOperation
	GetBlocks() []byte
}

type CommitEventType string

const (
	CommitEventTypeRepoCommit CommitEventType = "repo_commit"
)
