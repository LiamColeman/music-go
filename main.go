package main

import (
	"context"
	"errors"
	"log"
	"music-go/internal/handler"
	"music-go/internal/repository"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

type Artist struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

type ArtistWithAlbums struct {
	Artist
	Albums []Album
}

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

type AlbumWithSongs struct {
	Album
	Songs []Song
}

func getSongs(c *gin.Context) ([]Song, error) {

	query := `SELECT song.id, song.title, song.track_number, song.duration_seconds, album.name as album, artist.name as artist
				FROM song
				JOIN album ON song.album_id = album.id
				JOIN artist ON album.artist_id = artist.id
				ORDER BY song.title`
	rows, err := dbPool.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := []Song{}

	for rows.Next() {
		var song Song
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

func getSong(c *gin.Context, id string) (*Song, error) {

	query := `SELECT song.id, song.title, song.track_number, song.duration_seconds, album.name as album, artist.name as artist
				FROM song
				JOIN album ON song.album_id = album.id
				JOIN artist ON album.artist_id = artist.id
				WHERE song.id = $1`
	var song Song

	err := dbPool.QueryRow(c, query, id).Scan(&song.ID, &song.Title, &song.TrackNumber, &song.DurationSeconds, &song.AlbumName, &song.ArtistName)
	if err != nil {
		return nil, err
	}

	return &song, nil
}

func createSong(c *gin.Context, song CreateSong) (*SongResponse, error) {

	var songCreated SongResponse

	query := `INSERT INTO song (album_id, title, track_number, duration_seconds) VALUES ($1, $2, $3, $4) RETURNING id, title, track_number, duration_seconds`
	err := dbPool.QueryRow(c, query, song.AlbumID, song.Title, song.TrackNumber, song.DurationSeconds).Scan(&songCreated.ID, &songCreated.Title, &songCreated.TrackNumber, &songCreated.DurationSeconds)
	if err != nil {
		return nil, err
	}

	return &songCreated, nil

}

func updateSong(c *gin.Context, song UpdateSong, id string) (*SongResponse, error) {

	var updateSong SongResponse

	query := `UPDATE song SET title = $2, track_number = $3, duration_seconds = $4 WHERE id = $1 RETURNING id, title, track_number, duration_seconds`

	err := dbPool.QueryRow(c, query, id, song.Title, song.TrackNumber, song.DurationSeconds).Scan(&updateSong.ID, &updateSong.Title, &updateSong.TrackNumber, &updateSong.DurationSeconds)
	if err != nil {
		return nil, err
	}

	return &updateSong, nil
}

func patchSong(c *gin.Context, song PatchSong, id string) (*SongResponse, error) {

	var patchedSong SongResponse

	if song.Title != nil {
		queryName := `UPDATE song SET title = $2 WHERE id = $1 RETURNING id, title`
		err := dbPool.QueryRow(c, queryName, id, song.Title).Scan(&patchedSong.ID, &patchedSong.Title)
		if err != nil {
			return nil, err
		}
	}

	if song.TrackNumber != nil {
		queryReleaseYear := `UPDATE song SET track_number = $2 WHERE id = $1 RETURNING id, track_number`
		err := dbPool.QueryRow(c, queryReleaseYear, id, song.TrackNumber).Scan(&patchedSong.ID, &patchedSong.TrackNumber)
		if err != nil {
			return nil, err
		}
	}

	if song.DurationSeconds != nil {
		queryReleaseYear := `UPDATE song SET duration_seconds = $2 WHERE id = $1 RETURNING id, duration_seconds`
		err := dbPool.QueryRow(c, queryReleaseYear, id, song.DurationSeconds).Scan(&patchedSong.ID, &patchedSong.DurationSeconds)
		if err != nil {
			return nil, err
		}
	}

	return &patchedSong, nil
}

func deleteSong(c *gin.Context, id string) error {

	query := `DELETE FROM song where id = $1`

	_, err := dbPool.Query(c, query, id)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	var err error
	dbPool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}
	defer dbPool.Close()

	// TODO: Add other handlers and repos

	artistRepo := repository.NewArtistRepository(dbPool)
	artistHandler := handler.NewArtistHandler(artistRepo)

	albumRepo := repository.NewAlbumRepository(dbPool)
	albumHandler := handler.NewAlbumHandler(albumRepo)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/artists", artistHandler.GetAll)
	router.GET("/artists/:id", artistHandler.GetArtist)
	router.POST("/artists", artistHandler.CreateArtist)
	router.PUT("/artists/:id", artistHandler.UpdateArtist)
	router.PATCH("/artists/:id", artistHandler.PatchArtist)
	router.DELETE("/artists/:id", artistHandler.DeleteArtist)

	router.GET("/albums", albumHandler.GetAll)
	router.GET("/albums/:id", albumHandler.GetAlbum)
	router.POST("/albums", albumHandler.CreateAlbum)
	router.PUT("/albums/:id", albumHandler.UpdateAlbum)
	router.PATCH("/albums/:id", albumHandler.PatchAlbum)
	router.DELETE("/albums/:id", albumHandler.DeleteAlbum)

	router.GET("/songs", func(c *gin.Context) {
		songs, err := getSongs(c)

		if err != nil {
			log.Printf("Error fetching songs: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, songs)
	})

	router.POST("/songs", func(c *gin.Context) {
		var newSong CreateSong

		err := c.BindJSON(&newSong)
		if err != nil {
			return
		}

		songCreated, err := createSong(c, newSong)

		if err != nil {
			log.Printf("Error creating song %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		newUrl := "Location: /songs/" + strconv.Itoa(songCreated.ID)
		c.Header("location", newUrl)
		c.JSON(http.StatusCreated, songCreated)
	})

	router.GET("/songs/:id", func(c *gin.Context) {
		id := c.Param("id")
		song, err := getSong(c, id)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
				return
			}

			log.Printf("Error fetching song %s: %v", id, err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, song)
	})

	router.PUT("/songs/:id", func(c *gin.Context) {
		id := c.Param("id")
		var newSong UpdateSong

		err := c.BindJSON(&newSong)
		if err != nil {
			return
		}

		updatedSong, err := updateSong(c, newSong, id)

		if err != nil {
			log.Printf("Error updating song %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		newUrl := "Location: /songs/" + strconv.Itoa(updatedSong.ID)
		c.Header("location", newUrl)
		c.JSON(http.StatusCreated, updatedSong)
	})

	router.PATCH("/songs/:id", func(c *gin.Context) {
		id := c.Param("id")
		var newSong PatchSong

		err := c.BindJSON(&newSong)
		if err != nil {
			return
		}

		patchedSong, err := patchSong(c, newSong, id)

		if err != nil {
			log.Printf("Error patching song %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		newUrl := "Location: /songs/" + strconv.Itoa(patchedSong.ID)
		c.Header("location", newUrl)
		c.JSON(http.StatusCreated, patchedSong)
	})

	router.DELETE("/songs/:id", func(c *gin.Context) {
		id := c.Param("id")
		err := deleteSong(c, id)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
				return
			}

			log.Printf("Error deleting song %s: %v", id, err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusNoContent, "Deleted Song")
	})

	router.Run(":9000")
}
