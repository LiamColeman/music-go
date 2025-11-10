package handler

import (
	"errors"
	"log"
	"music-go/internal/model"
	"music-go/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type SongHandler struct {
	songRepo *repository.SongRepository
}

func NewSongHandler(songRepo *repository.SongRepository) *SongHandler {
	return &SongHandler{
		songRepo: songRepo,
	}
}

func (h *SongHandler) GetAll(c *gin.Context) {
	songs, err := h.songRepo.GetSongs(c.Request.Context())

	if err != nil {
		log.Printf("Error fetching songs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

func (h *SongHandler) CreateSong(c *gin.Context) {
	var newSong model.CreateSong

	err := c.BindJSON(&newSong)
	if err != nil {
		return
	}

	songCreated, err := h.songRepo.CreateSong(c.Request.Context(), newSong)

	if err != nil {
		log.Printf("Error creating song %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newUrl := "Location: /songs/" + strconv.Itoa(songCreated.ID)
	c.Header("location", newUrl)
	c.JSON(http.StatusCreated, songCreated)
}

func (h *SongHandler) GetSong(c *gin.Context) {
	id := c.Param("id")
	song, err := h.songRepo.GetSong(c.Request.Context(), id)

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
}

func (h *SongHandler) UpdateSong(c *gin.Context) {
	id := c.Param("id")
	var newSong model.UpdateSong

	err := c.BindJSON(&newSong)
	if err != nil {
		return
	}

	updatedSong, err := h.songRepo.UpdateSong(c.Request.Context(), newSong, id)

	if err != nil {
		log.Printf("Error updating song %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newUrl := "Location: /songs/" + strconv.Itoa(updatedSong.ID)
	c.Header("location", newUrl)
	c.JSON(http.StatusCreated, updatedSong)
}

func (h *SongHandler) PatchSong(c *gin.Context) {
	id := c.Param("id")
	var newSong model.PatchSong

	err := c.BindJSON(&newSong)
	if err != nil {
		return
	}

	patchedSong, err := h.songRepo.PatchSong(c.Request.Context(), newSong, id)

	if err != nil {
		log.Printf("Error patching song %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newUrl := "Location: /songs/" + strconv.Itoa(patchedSong.ID)
	c.Header("location", newUrl)
	c.JSON(http.StatusCreated, patchedSong)
}

func (h *SongHandler) DeleteSong(c *gin.Context) {
	id := c.Param("id")
	err := h.songRepo.DeleteSong(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}

		log.Printf("Error deleting song %s: %v", id, err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
