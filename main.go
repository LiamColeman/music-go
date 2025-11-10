package main

import (
	"context"
	"log"
	"music-go/internal/handler"
	"music-go/internal/repository"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

func main() {

	var err error
	dbPool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}
	defer dbPool.Close()

	artistRepo := repository.NewArtistRepository(dbPool)
	artistHandler := handler.NewArtistHandler(artistRepo)

	albumRepo := repository.NewAlbumRepository(dbPool)
	albumHandler := handler.NewAlbumHandler(albumRepo)

	songRepo := repository.NewSongRepository(dbPool)
	songHandler := handler.NewSongHandler(songRepo)

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

	router.GET("/songs", songHandler.GetAll)
	router.GET("/songs/:id", songHandler.GetSong)
	router.POST("/songs", songHandler.CreateSong)
	router.PUT("/songs/:id", songHandler.UpdateSong)
	router.PATCH("/songs/:id", songHandler.PatchSong)
	router.DELETE("/songs/:id", songHandler.DeleteSong)

	router.Run(":9000")

}
