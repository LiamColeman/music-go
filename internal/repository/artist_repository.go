package repository

import (
	"context"
	"music-go/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ArtistRepository struct {
	dbPool *pgxpool.Pool
}

func NewArtistRepository(dbPool *pgxpool.Pool) *ArtistRepository {
	return &ArtistRepository{
		dbPool: dbPool,
	}
}

func (r *ArtistRepository) GetArtists(ctx context.Context) ([]model.Artist, error) {
	query := `SELECT id, name, description FROM artist`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	artists := []model.Artist{}

	for rows.Next() {
		var artist model.Artist
		if err := rows.Scan(&artist.ID, &artist.Name, &artist.Description); err != nil {
			return nil, err
		}

		artists = append(artists, artist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return artists, nil
}

func (r *ArtistRepository) getArtist(ctx context.Context, id string) (*model.ArtistWithAlbums, error) {
	var artist model.ArtistWithAlbums
	query := `SELECT id, name, description FROM artist WHERE id = $1`

	err := r.dbPool.QueryRow(ctx, query, id).Scan(&artist.ID, &artist.Name, &artist.Description)
	if err != nil {
		return nil, err
	}

	artist.Albums, err = r.getAlbumsForArtist(ctx, artist.ID)
	if err != nil {
		return nil, err
	}

	return &artist, nil

}

func (r *ArtistRepository) getAlbumsForArtist(ctx context.Context, artistID int) ([]model.Album, error) {

	query := `SELECT id, name, release_year FROM album WHERE artist_id = $1 ORDER BY release_year DESC`
	rows, err := r.dbPool.Query(ctx, query, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	albums := []model.Album{}

	for rows.Next() {
		var album model.Album
		if err := rows.Scan(&album.ID, &album.Name, &album.ReleaseYear); err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}
