package repository

import (
	"context"

	"github.com/liamcoleman/music-go/internal/model"

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

func (r *ArtistRepository) GetArtist(ctx context.Context, id string) (*model.ArtistWithAlbums, error) {
	var artist model.ArtistWithAlbums
	query := `SELECT id, name, description FROM artist WHERE id = $1`

	err := r.dbPool.QueryRow(ctx, query, id).Scan(&artist.ID, &artist.Name, &artist.Description)
	if err != nil {
		return nil, err
	}

	artist.Albums, err = r.GetAlbumsForArtist(ctx, artist.ID)
	if err != nil {
		return nil, err
	}

	return &artist, nil

}

func (r *ArtistRepository) GetAlbumsForArtist(ctx context.Context, artistID int) ([]model.Album, error) {

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

func (r *ArtistRepository) CreateArtist(ctx context.Context, artist model.CreateArtist) (*model.Artist, error) {

	var createdArtist model.Artist
	query := `INSERT INTO artist (name, description) VALUES ($1, $2) RETURNING id, name, description`
	err := r.dbPool.QueryRow(ctx, query, artist.Name, artist.Description).Scan(&createdArtist.ID, &createdArtist.Name, &createdArtist.Description)
	if err != nil {
		return nil, err
	}

	return &createdArtist, nil

}

func (r *ArtistRepository) UpdateArtist(ctx context.Context, artist model.UpdateArtist, id string) (*model.Artist, error) {
	var updatedArtist model.Artist

	query := `UPDATE artist SET name = $2, description = $3 WHERE id = $1 RETURNING id, name, description`

	err := r.dbPool.QueryRow(ctx, query, id, artist.Name, artist.Description).Scan(&updatedArtist.ID, &updatedArtist.Name, &updatedArtist.Description)
	if err != nil {
		return nil, err
	}

	return &updatedArtist, nil

}

func (r *ArtistRepository) PatchArtist(ctx context.Context, artist model.PatchArtist, id string) (*model.Artist, error) {

	var patchedArtist model.Artist

	if artist.Name != nil {
		queryName := `UPDATE artist SET name = $2 WHERE id = $1`
		_, err := r.dbPool.Exec(ctx, queryName, id, artist.Name)
		if err != nil {
			return nil, err
		}
	}

	if artist.Description != nil {
		queryDescription := `UPDATE artist SET description = $2 WHERE id = $1`
		_, err := r.dbPool.Exec(ctx, queryDescription, id, artist.Description)
		if err != nil {
			return nil, err
		}
	}

	query := `SELECT id, name, description FROM artist WHERE id = $1`

	err := r.dbPool.QueryRow(ctx, query, id).Scan(&patchedArtist.ID, &patchedArtist.Name, &patchedArtist.Description)
	if err != nil {
		return nil, err
	}

	return &patchedArtist, nil
}

func (r *ArtistRepository) DeleteArtist(ctx context.Context, id string) error {

	query := `DELETE FROM artist where id = $1`

	_, err := r.dbPool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
