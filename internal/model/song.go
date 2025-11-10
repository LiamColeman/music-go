package model

type Song struct {
	ID              int    `json:"id"`
	ArtistName      string `json:"artist,omitempty"`
	AlbumName       string `json:"album,omitempty"`
	Title           string `json:"title"`
	TrackNumber     int    `json:"track_number"`
	DurationSeconds int    `json:"duration_seconds"`
}

type SongResponse struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	TrackNumber     int    `json:"track_number"`
	DurationSeconds int    `json:"duration_seconds"`
}

type CreateSong struct {
	AlbumID         int    `json:"album_id"`
	Title           string `json:"title"`
	TrackNumber     int    `json:"track_number"`
	DurationSeconds int    `json:"duration_seconds"`
}

type UpdateSong struct {
	Title           string `json:"title"`
	TrackNumber     int    `json:"track_number"`
	DurationSeconds int    `json:"duration_seconds"`
}

type PatchSong struct {
	Title           *string `json:"title"`
	TrackNumber     *int    `json:"track_number"`
	DurationSeconds *int    `json:"duration_seconds"`
}
