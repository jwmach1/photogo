package data

import "time"

type ListRequest struct {
	PageToken string `json:"pageToken,omitempty"`
	PageSize  string `json:"pageSize,omitempty"`
}
type MediaResponse struct {
	MediaItems    []*MediaItem `json:"mediaItems"`
	NextPageToken string       `json:"nextPageToken"`
}
type MediaItem struct {
	ID       string        `json:"id"`
	Filename string        `json:"filename"`
	BaseUrl  string        `json:"baseUrl"`
	MimeType string        `json:"mimeType"`
	Metadata MediaMetadata `json:"mediaMetadata"`
}

type MediaMetadata struct {
	CreationTime time.Time `json:"creationTime"`
}
