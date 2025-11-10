package model

type Album struct {
	ID          int    `json:"id"`
	ArtistName  string `json:"artist,omitempty"`
	Name        string `json:"name"`
	ReleaseYear int    `json:"release_year"`
}

type AlbumResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ReleaseYear int    `json:"release_year"`
}

type CreateAlbum struct {
	ArtistID    int    `json:"artist_id"`
	Name        string `json:"name"`
	ReleaseYear int    `json:"release_year"`
}

type UpdateAlbum struct {
	Name        string `json:"name"`
	ReleaseYear int    `json:"release_year"`
}

type PatchAlbum struct {
	Name        *string `json:"name"`
	ReleaseYear *int    `json:"release_year"`
}

type AlbumWithSongs struct {
	Album
	Songs []Song
}
