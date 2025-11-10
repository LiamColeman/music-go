package repository

import (
	"context"

	"github.com/liamcoleman/music-go/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AlbumRepository struct {
	dbPool *pgxpool.Pool
}

func NewAlbumRepository(dbPool *pgxpool.Pool) *AlbumRepository {
	return &AlbumRepository{
		dbPool: dbPool,
	}
}

func (r *AlbumRepository) GetAlbums(ctx context.Context) ([]model.Album, error) {

	query := `SELECT album.id, album.name, album.release_year, artist.name as artist 
			FROM album 
			JOIN artist ON album.artist_id = artist.id 
			ORDER BY artist.name, album.name`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	albums := []model.Album{}

	for rows.Next() {
		var album model.Album
		if err := rows.Scan(&album.ID, &album.Name, &album.ReleaseYear, &album.ArtistName); err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

func (r *AlbumRepository) GetAlbum(ctx context.Context, id string) (*model.AlbumWithSongs, error) {

	var album model.AlbumWithSongs

	query := `SELECT album.id, album.name, album.release_year, artist.name as artist 
		FROM album 
		JOIN artist ON album.artist_id = artist.id 
		WHERE album.id = $1`

	err := r.dbPool.QueryRow(ctx, query, id).Scan(&album.ID, &album.Name, &album.ReleaseYear, &album.ArtistName)
	if err != nil {
		return nil, err
	}

	album.Songs, err = r.GetSongsForAlbum(ctx, album.ID)
	if err != nil {
		return nil, err
	}

	return &album, nil

}

func (r *AlbumRepository) CreateAlbum(ctx context.Context, album model.CreateAlbum) (*model.AlbumResponse, error) {

	var albumCreated model.AlbumResponse

	query := `INSERT INTO album (artist_id, name, release_year) VALUES ($1, $2, $3) RETURNING id, name, release_year`
	err := r.dbPool.QueryRow(ctx, query, album.ArtistID, album.Name, album.ReleaseYear).Scan(&albumCreated.ID, &albumCreated.Name, &albumCreated.ReleaseYear)
	if err != nil {
		return nil, err
	}

	return &albumCreated, nil

}

func (r *AlbumRepository) UpdateAlbum(ctx context.Context, album model.UpdateAlbum, id string) (*model.AlbumResponse, error) {
	var updatedAlbum model.AlbumResponse

	query := `UPDATE album SET name = $2, release_year = $3 WHERE id = $1 RETURNING id, name, release_year`

	err := r.dbPool.QueryRow(ctx, query, id, album.Name, album.ReleaseYear).Scan(&updatedAlbum.ID, &updatedAlbum.Name, &updatedAlbum.ReleaseYear)
	if err != nil {
		return nil, err
	}

	return &updatedAlbum, nil
}

func (r *AlbumRepository) PatchAlbum(ctx context.Context, album model.PatchAlbum, id string) (*model.AlbumResponse, error) {

	var patchedAlbum model.AlbumResponse

	if album.Name != nil {
		queryName := `UPDATE album SET name = $2 WHERE id = $1`
		_, err := r.dbPool.Exec(ctx, queryName, id, album.Name)
		if err != nil {
			return nil, err
		}
	}

	if album.ReleaseYear != nil {
		queryReleaseYear := `UPDATE album SET release_year = $2 WHERE id = $1`
		_, err := r.dbPool.Exec(ctx, queryReleaseYear, id, album.ReleaseYear)
		if err != nil {
			return nil, err
		}
	}

	query := `SELECT id, name, release_year FROM album WHERE id = $1`

	err := r.dbPool.QueryRow(ctx, query, id).Scan(&patchedAlbum.ID, &patchedAlbum.Name, &patchedAlbum.ReleaseYear)
	if err != nil {
		return nil, err
	}

	return &patchedAlbum, nil
}

func (r *AlbumRepository) GetSongsForAlbum(ctx context.Context, albumID int) ([]model.Song, error) {

	query := `SELECT id, title, track_number, duration_seconds FROM song WHERE album_id = $1 ORDER BY track_number`
	rows, err := r.dbPool.Query(ctx, query, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := []model.Song{}

	for rows.Next() {
		var song model.Song
		if err := rows.Scan(&song.ID, &song.Title, &song.TrackNumber, &song.DurationSeconds); err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (r *AlbumRepository) DeleteAlbum(ctx context.Context, id string) error {

	query := `DELETE FROM album where id = $1`

	_, err := r.dbPool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
