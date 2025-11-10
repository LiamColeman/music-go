package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/liamcoleman/music-go/internal/repository"

	"github.com/liamcoleman/music-go/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AlbumHandler struct {
	albumRepo *repository.AlbumRepository
}

func NewAlbumHandler(albumRepo *repository.AlbumRepository) *AlbumHandler {
	return &AlbumHandler{
		albumRepo: albumRepo,
	}
}

func (h *AlbumHandler) GetAll(c *gin.Context) {
	albums, err := h.albumRepo.GetAlbums(c.Request.Context())

	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, albums)
}

func (h *AlbumHandler) GetAlbum(c *gin.Context) {
	id := c.Param("id")
	album, err := h.albumRepo.GetAlbum(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}

		log.Printf("Error fetching album %s: %v", id, err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, album)
}

func (h *AlbumHandler) CreateAlbum(c *gin.Context) {
	var newAlbum model.CreateAlbum

	err := c.BindJSON(&newAlbum)
	if err != nil {
		return
	}

	albumCreated, err := h.albumRepo.CreateAlbum(c.Request.Context(), newAlbum)

	if err != nil {
		log.Printf("Error creating album %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newUrl := "Location: /albums/" + strconv.Itoa(albumCreated.ID)
	c.Header("location", newUrl)
	c.JSON(http.StatusCreated, albumCreated)
}

func (h *AlbumHandler) UpdateAlbum(c *gin.Context) {
	id := c.Param("id")
	var newAlbum model.UpdateAlbum

	err := c.BindJSON(&newAlbum)
	if err != nil {
		return
	}

	updatedAlbum, err := h.albumRepo.UpdateAlbum(c.Request.Context(), newAlbum, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}

		log.Printf("Error updating album %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newUrl := "Location: /artists/" + strconv.Itoa(updatedAlbum.ID)
	c.Header("location", newUrl)
	c.JSON(http.StatusOK, updatedAlbum)
}

func (h *AlbumHandler) PatchAlbum(c *gin.Context) {
	id := c.Param("id")
	var newAlbum model.PatchAlbum

	err := c.BindJSON(&newAlbum)
	if err != nil {
		return
	}

	patchedAlbum, err := h.albumRepo.PatchAlbum(c.Request.Context(), newAlbum, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}

		log.Printf("Error patching album %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newUrl := "Location: /artists/" + strconv.Itoa(patchedAlbum.ID)
	c.Header("location", newUrl)
	c.JSON(http.StatusOK, patchedAlbum)
}

func (h *AlbumHandler) DeleteAlbum(c *gin.Context) {
	id := c.Param("id")
	err := h.albumRepo.DeleteAlbum(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}

		log.Printf("Error deleting album %s: %v", id, err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
