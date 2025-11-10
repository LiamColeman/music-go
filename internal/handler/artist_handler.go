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

func (h *ArtistHandler) GetArtist(c *gin.Context) {
	id := c.Param("id")
	artist, err := h.artistRepo.GetArtist(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
			return
		}

		log.Printf("Error fetching artist %s: %v", id, err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, artist)
}

func (h *ArtistHandler) CreateArtist(c *gin.Context) {
	var newArtist model.CreateArtist

	err := c.BindJSON(&newArtist)
	if err != nil {
		return
	}

	createdArtist, err := h.artistRepo.CreateArtist(c.Request.Context(), newArtist)

	if err != nil {
		log.Printf("Error creating artist %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newUrl := "Location: /artists/" + strconv.Itoa(createdArtist.ID)
	c.Header("location", newUrl)
	c.JSON(http.StatusCreated, createdArtist)
}

func (h *ArtistHandler) UpdateArtist(c *gin.Context) {
	id := c.Param("id")
	var newArtist model.UpdateArtist

	err := c.BindJSON(&newArtist)
	if err != nil {
		return
	}

	updatedArtist, err := h.artistRepo.UpdateArtist(c.Request.Context(), newArtist, id)

	if err != nil {
		log.Printf("Error updating artist %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newUrl := "Location: /artists/" + strconv.Itoa(updatedArtist.ID)
	c.Header("location", newUrl)
	c.JSON(http.StatusOK, updatedArtist)
}

func (h *ArtistHandler) PatchArtist(c *gin.Context) {
	id := c.Param("id")

	var newArtist model.PatchArtist

	err := c.BindJSON(&newArtist)
	if err != nil {
		return
	}

	patchedArtist, err := h.artistRepo.PatchArtist(c.Request.Context(), newArtist, id)

	if err != nil {
		log.Printf("Error patching artist %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newUrl := "Location: /artists/" + strconv.Itoa(patchedArtist.ID)
	c.Header("location", newUrl)
	c.JSON(http.StatusOK, patchedArtist)
}

func (h *ArtistHandler) DeleteArtist(c *gin.Context) {
	id := c.Param("id")
	err := h.artistRepo.DeleteArtist(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
			return
		}

		log.Printf("Error deleting artist %s: %v", id, err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, "Deleted Artist")
}
