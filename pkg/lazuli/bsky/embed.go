package bsky

type EmbedRecord struct {
	LexiconTypeID string `json:"$type"`
	Record        Record `json:"record"`
}

type BlobRef struct {
	Link string `json:"$link"`
}

type BlobRecord struct {
	LexiconTypeID string  `json:"$type"`
	Ref           BlobRef `json:"ref"`
	MimeType      string  `json:"mimeType"`
	Size          int     `json:"size"`
}

type ImageAspectRatio struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type ImageRecord struct {
	Alt         string           `json:"alt"`
	AspectRatio ImageAspectRatio `json:"aspectRatio"`
	Image       BlobRecord       `json:"image"`
}

type EmbedImageRecord struct {
	LexiconTypeID string        `json:"$type"`
	Images        []ImageRecord `json:"images"`
}
