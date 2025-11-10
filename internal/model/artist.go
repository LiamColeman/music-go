package model

type Artist struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ArtistWithAlbums struct {
	Artist
	Albums []Album
}

type CreateArtist struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateArtist struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type PatchArtist struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
