package bsky

type Reply struct {
	Parent RepoStrongRef `json:"parent"`
	Root   RepoStrongRef `json:"root"`
}
