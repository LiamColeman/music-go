package handler

import (
	"log"
	"music-go/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArtistHandler struct {
	artistRepo *repository.ArtistRepository
}

func NewArtistHandler(artistRepo *repository.ArtistRepository) *ArtistHandler {
	return &ArtistHandler{
		artistRepo: artistRepo,
	}
}

func (h *ArtistHandler) GetAll(c *gin.Context) {
	artists, err := h.artistRepo.GetArtists(c.Request.Context())

	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, artists)
}
