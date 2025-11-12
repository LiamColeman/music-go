package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/liamcoleman/music-go/internal/model"
	"github.com/liamcoleman/music-go/internal/repository"
	"github.com/stretchr/testify/assert"
)

// setupTestDBAlbum creates a connection pool to your test database
// NOTE: Start test database with: docker compose -f docker-compose.test.yml up -d
func setupTestDBAlbum(t *testing.T) *pgxpool.Pool {
	// Use test database URL (separate DB on port 5433)
	databaseURL := "postgresql://postgres:gizzard@localhost:5433/albums_test"

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Fatalf("Unable to create connection pool: %v", err)
	}

	return pool
}

// cleanupAlbums removes test data from the database
func cleanupAlbums(t *testing.T, pool *pgxpool.Pool, name string) {
	_, err := pool.Exec(context.Background(),
		"DELETE FROM album WHERE name = $1", name)
	if err != nil {
		t.Logf("Warning: cleanup failed: %v", err)
	}
}

func TestAlbums(t *testing.T) {
	// Setup
	pool := setupTestDBAlbum(t)
	defer pool.Close()

	albumRepo := repository.NewAlbumRepository(pool)
	handler := NewAlbumHandler(albumRepo)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	TestArtistID := 1
	testAlbumName := "Test album for Integration"
	testAlbumReleaseYear := 2019
	defer cleanupAlbums(t, pool, testAlbumName)

	// Store this so we can use in tests
	var createdAlbum model.Album

	t.Run("CreateAlbum", func(t *testing.T) {
		// Create request body
		createAlbum := model.CreateAlbum{
			ArtistID:    TestArtistID,
			Name:        testAlbumName,
			ReleaseYear: testAlbumReleaseYear,
		}

		body, _ := json.Marshal(createAlbum)

		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/albums", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call handler
		handler.CreateAlbum(c)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &createdAlbum)
		assert.NoError(t, err)
		assert.Equal(t, testAlbumName, createdAlbum.Name)
		assert.Equal(t, testAlbumReleaseYear, createdAlbum.ReleaseYear)
		assert.NotZero(t, createdAlbum.ID)

		// Check Location header
		location := w.Header().Get("Location")
		locationUrl := "/albums/" + strconv.Itoa(createdAlbum.ID)
		assert.Contains(t, location, locationUrl)
	})

	t.Run("GetAlbum", func(t *testing.T) {
		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var getAlbumUrl = "/albums/" + strconv.Itoa(createdAlbum.ID)
		c.Request = httptest.NewRequest("GET", getAlbumUrl, nil)

		// Set the URL parameter that the handler expects
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdAlbum.ID)}}

		// Call handler
		handler.GetAlbum(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var album model.Album
		err := json.Unmarshal(w.Body.Bytes(), &album)
		assert.NoError(t, err)
		assert.Equal(t, createdAlbum.Name, album.Name)
		assert.Equal(t, createdAlbum.ID, album.ID)
	})

	t.Run("UpdateAlbum", func(t *testing.T) {

		updatedAlbumName := "Update album for Integration"
		updatedAlbumReleaseYear := 2018

		// Store this so we can use in tests
		var updatedAlbum model.Album

		// Create request body
		updateAlbum := model.UpdateAlbum{
			Name:        updatedAlbumName,
			ReleaseYear: updatedAlbumReleaseYear,
		}

		body, _ := json.Marshal(updateAlbum)

		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var updateAlbumUrl = "/albums/" + strconv.Itoa(createdAlbum.ID)
		c.Request = httptest.NewRequest("PUT", updateAlbumUrl, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdAlbum.ID)}}

		// Call handler
		handler.UpdateAlbum(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &updatedAlbum)
		assert.NoError(t, err)
		assert.Equal(t, updatedAlbumName, updatedAlbum.Name)
		assert.Equal(t, updatedAlbumReleaseYear, updatedAlbum.ReleaseYear)
		assert.NotZero(t, updatedAlbum.ID)

		// Check Location header
		location := w.Header().Get("Location")
		locationUrl := "/albums/" + strconv.Itoa(updatedAlbum.ID)
		assert.Contains(t, location, locationUrl)
	})

	t.Run("PatchAlbum", func(t *testing.T) {

		patchedAlbumName := "Patched album for Integration"
		patchedAlbumReleaseYear := 2017

		// Store this so we can use in tests
		var patchedAlbum model.Album

		// Create request body
		patchAlbum := model.PatchAlbum{
			Name:        &patchedAlbumName,
			ReleaseYear: &patchedAlbumReleaseYear,
		}

		body, _ := json.Marshal(patchAlbum)

		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var patchAlbumUrl = "/albums/" + strconv.Itoa(createdAlbum.ID)
		c.Request = httptest.NewRequest("Patch", patchAlbumUrl, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdAlbum.ID)}}

		// Call handler
		handler.PatchAlbum(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &patchedAlbum)
		assert.NoError(t, err)
		assert.Equal(t, patchedAlbumName, patchedAlbum.Name)
		assert.Equal(t, patchedAlbumReleaseYear, patchedAlbum.ReleaseYear)
		assert.NotZero(t, patchedAlbum.ID)

		// Check Location header
		location := w.Header().Get("Location")
		locationUrl := "/albums/" + strconv.Itoa(patchedAlbum.ID)
		assert.Contains(t, location, locationUrl)
	})

}
