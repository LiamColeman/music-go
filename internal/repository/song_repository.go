package repository

import (
	"context"
	"music-go/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SongRepository struct {
	dbPool *pgxpool.Pool
}

func NewSongRepository(dbPool *pgxpool.Pool) *SongRepository {
	return &SongRepository{
		dbPool: dbPool,
	}
}

func (r *SongRepository) GetSongs(ctx context.Context) ([]model.Song, error) {

	query := `SELECT song.id, song.title, song.track_number, song.duration_seconds, album.name as album, artist.name as artist
				FROM song
				JOIN album ON song.album_id = album.id
				JOIN artist ON album.artist_id = artist.id
				ORDER BY song.title`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := []model.Song{}

	for rows.Next() {
		var song model.Song
		if err := rows.Scan(&song.ID, &song.Title, &song.TrackNumber, &song.DurationSeconds, &song.AlbumName, &song.ArtistName); err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (r *SongRepository) GetSong(ctx context.Context, id string) (*model.Song, error) {

	query := `SELECT song.id, song.title, song.track_number, song.duration_seconds, album.name as album, artist.name as artist
				FROM song
				JOIN album ON song.album_id = album.id
				JOIN artist ON album.artist_id = artist.id
				WHERE song.id = $1`
	var song model.Song

	err := r.dbPool.QueryRow(ctx, query, id).Scan(&song.ID, &song.Title, &song.TrackNumber, &song.DurationSeconds, &song.AlbumName, &song.ArtistName)
	if err != nil {
		return nil, err
	}

	return &song, nil
}

func (r *SongRepository) CreateSong(ctx context.Context, song model.CreateSong) (*model.SongResponse, error) {

	var songCreated model.SongResponse

	query := `INSERT INTO song (album_id, title, track_number, duration_seconds) VALUES ($1, $2, $3, $4) RETURNING id, title, track_number, duration_seconds`
	err := r.dbPool.QueryRow(ctx, query, song.AlbumID, song.Title, song.TrackNumber, song.DurationSeconds).Scan(&songCreated.ID, &songCreated.Title, &songCreated.TrackNumber, &songCreated.DurationSeconds)
	if err != nil {
		return nil, err
	}

	return &songCreated, nil

}

func (r *SongRepository) UpdateSong(ctx context.Context, song model.UpdateSong, id string) (*model.SongResponse, error) {

	var updateSong model.SongResponse

	query := `UPDATE song SET title = $2, track_number = $3, duration_seconds = $4 WHERE id = $1 RETURNING id, title, track_number, duration_seconds`

	err := r.dbPool.QueryRow(ctx, query, id, song.Title, song.TrackNumber, song.DurationSeconds).Scan(&updateSong.ID, &updateSong.Title, &updateSong.TrackNumber, &updateSong.DurationSeconds)
	if err != nil {
		return nil, err
	}

	return &updateSong, nil
}

func (r *SongRepository) PatchSong(ctx context.Context, song model.PatchSong, id string) (*model.SongResponse, error) {

	var patchedSong model.SongResponse

	if song.Title != nil {
		queryName := `UPDATE song SET title = $2 WHERE id = $1 RETURNING id, title`
		err := r.dbPool.QueryRow(ctx, queryName, id, song.Title).Scan(&patchedSong.ID, &patchedSong.Title)
		if err != nil {
			return nil, err
		}
	}

	if song.TrackNumber != nil {
		queryReleaseYear := `UPDATE song SET track_number = $2 WHERE id = $1 RETURNING id, track_number`
		err := r.dbPool.QueryRow(ctx, queryReleaseYear, id, song.TrackNumber).Scan(&patchedSong.ID, &patchedSong.TrackNumber)
		if err != nil {
			return nil, err
		}
	}

	if song.DurationSeconds != nil {
		queryReleaseYear := `UPDATE song SET duration_seconds = $2 WHERE id = $1 RETURNING id, duration_seconds`
		err := r.dbPool.QueryRow(ctx, queryReleaseYear, id, song.DurationSeconds).Scan(&patchedSong.ID, &patchedSong.DurationSeconds)
		if err != nil {
			return nil, err
		}
	}

	return &patchedSong, nil
}

func (r *SongRepository) DeleteSong(ctx context.Context, id string) error {

	query := `DELETE FROM song where id = $1`

	_, err := r.dbPool.Query(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
