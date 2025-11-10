package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/liamcoleman/music-go/internal/repository"

	"github.com/liamcoleman/music-go/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
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

	router.Run()

}
