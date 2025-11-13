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

// setupTestDBSong creates a connection pool to your test database
// NOTE: Start test database with: docker compose -f docker-compose.test.yml up -d
func setupTestDBSong(t *testing.T) *pgxpool.Pool {
	// Use test database URL (separate DB on port 5433)
	databaseURL := "postgresql://postgres:gizzard@localhost:5433/albums_test"

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Fatalf("Unable to create connection pool: %v", err)
	}

	return pool
}

// cleanupSongs removes test data from the database
func cleanupSongs(t *testing.T, pool *pgxpool.Pool, name string) {
	_, err := pool.Exec(context.Background(),
		"DELETE FROM song WHERE title = $1", name)
	if err != nil {
		t.Logf("Warning: cleanup failed: %v", err)
	}
}

func TestSongs(t *testing.T) {
	// Setup
	pool := setupTestDBSong(t)
	defer pool.Close()

	songRepo := repository.NewSongRepository(pool)
	handler := NewSongHandler(songRepo)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	testSongTitle := "Test Song for Integration"
	testAlbumID := 1
	testTrackNumber := 1
	testDurationSeconds := 60

	defer cleanupSongs(t, pool, testSongTitle)

	// Store this so we can use in tests
	var createdSong model.Song

	t.Run("CreateSong", func(t *testing.T) {
		// Create request body
		createSong := model.CreateSong{
			AlbumID:         testAlbumID,
			Title:           testSongTitle,
			TrackNumber:     testTrackNumber,
			DurationSeconds: testDurationSeconds,
		}

		body, _ := json.Marshal(createSong)

		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/songs", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call handler
		handler.CreateSong(c)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &createdSong)
		assert.NoError(t, err)
		assert.Equal(t, testSongTitle, createdSong.Title)
		assert.Equal(t, testTrackNumber, createdSong.TrackNumber)
		assert.Equal(t, testDurationSeconds, createdSong.DurationSeconds)
		assert.NotZero(t, createdSong.ID)

		// Check Location header
		location := w.Header().Get("Location")
		locationUrl := "/songs/" + strconv.Itoa(createdSong.ID)
		assert.Contains(t, location, locationUrl)
	})

	t.Run("GetSong", func(t *testing.T) {
		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var getSongUrl = "/songs/" + strconv.Itoa(createdSong.ID)
		c.Request = httptest.NewRequest("GET", getSongUrl, nil)

		// Set the URL parameter that the handler expects
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdSong.ID)}}

		// Call handler
		handler.GetSong(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var song model.Song
		err := json.Unmarshal(w.Body.Bytes(), &song)
		assert.NoError(t, err)
		assert.Equal(t, testSongTitle, song.Title)
		assert.Equal(t, testTrackNumber, song.TrackNumber)
		assert.Equal(t, testDurationSeconds, song.DurationSeconds)
		assert.Equal(t, createdSong.ID, song.ID)
	})

	t.Run("GetAllSongs", func(t *testing.T) {
		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/songs", nil)

		// Call handler
		handler.GetAll(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var songs []model.Song
		err := json.Unmarshal(w.Body.Bytes(), &songs)
		assert.NoError(t, err)
		assert.NotEmpty(t, songs)

		// Verify our test song is in the list
		found := false
		for _, song := range songs {
			if song.ID == createdSong.ID {
				found = true
				assert.Equal(t, testSongTitle, song.Title)
				assert.Equal(t, testTrackNumber, song.TrackNumber)
				assert.Equal(t, testDurationSeconds, song.DurationSeconds)
				assert.Equal(t, createdSong.ID, song.ID)
				break
			}
		}
		assert.True(t, found, testSongTitle)
	})

	t.Run("UpdateSong", func(t *testing.T) {

		updatedSongTitle := "Update Song for Integration"
		updatedSongTrackNumber := 2
		updatedDurationSeconds := 61

		// Store this so we can use in tests
		var updatedSong model.Song

		// Create request body
		updateSong := model.UpdateSong{
			Title:           updatedSongTitle,
			TrackNumber:     updatedSongTrackNumber,
			DurationSeconds: updatedDurationSeconds,
		}

		body, _ := json.Marshal(updateSong)

		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var updateSongUrl = "/songs/" + strconv.Itoa(createdSong.ID)
		c.Request = httptest.NewRequest("PUT", updateSongUrl, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdSong.ID)}}

		// Call handler
		handler.UpdateSong(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &updatedSong)
		assert.NoError(t, err)
		assert.Equal(t, updatedSongTitle, updatedSong.Title)
		assert.Equal(t, updatedSongTrackNumber, updatedSong.TrackNumber)
		assert.Equal(t, updatedDurationSeconds, updatedSong.DurationSeconds)
		assert.Equal(t, createdSong.ID, updatedSong.ID)

		// Check Location header
		location := w.Header().Get("Location")
		locationUrl := "/songs/" + strconv.Itoa(updatedSong.ID)
		assert.Contains(t, location, locationUrl)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		// Verify the update happened by doing a get request
		var getSongUrl = "/songs/" + strconv.Itoa(updatedSong.ID)
		c.Request = httptest.NewRequest("GET", getSongUrl, nil)

		// Set the URL parameter that the handler expects
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(updatedSong.ID)}}

		// Call handler
		handler.GetSong(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var song model.Song
		err = json.Unmarshal(w.Body.Bytes(), &song)
		assert.NoError(t, err)
		assert.Equal(t, updatedSongTitle, song.Title)
		assert.Equal(t, updatedSongTrackNumber, song.TrackNumber)
		assert.Equal(t, updatedDurationSeconds, song.DurationSeconds)
		assert.Equal(t, createdSong.ID, song.ID)
	})

	t.Run("PatchSong", func(t *testing.T) {

		patchedSongTitle := "Patch Song for Integration"
		patchedSongTrackNumber := 2
		patchedDurationSeconds := 61

		// Store this so we can use in tests
		var patchedSong model.Song

		// Create request body
		patchSong := model.PatchSong{
			Title:           &patchedSongTitle,
			TrackNumber:     &patchedSongTrackNumber,
			DurationSeconds: &patchedDurationSeconds,
		}

		body, _ := json.Marshal(patchSong)

		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var updateSongUrl = "/songs/" + strconv.Itoa(createdSong.ID)
		c.Request = httptest.NewRequest("PATCH", updateSongUrl, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdSong.ID)}}

		// Call handler
		handler.PatchSong(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &patchedSong)
		assert.NoError(t, err)
		assert.Equal(t, patchedSongTitle, patchedSong.Title)
		assert.Equal(t, patchedSongTrackNumber, patchedSong.TrackNumber)
		assert.Equal(t, patchedDurationSeconds, patchedSong.DurationSeconds)
		assert.Equal(t, createdSong.ID, patchedSong.ID)

		// Check Location header
		location := w.Header().Get("Location")
		locationUrl := "/songs/" + strconv.Itoa(patchedSong.ID)
		assert.Contains(t, location, locationUrl)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		// Verify the update happened by doing a get request
		var getSongUrl = "/songs/" + strconv.Itoa(patchedSong.ID)
		c.Request = httptest.NewRequest("GET", getSongUrl, nil)

		// Set the URL parameter that the handler expects
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(patchedSong.ID)}}

		// Call handler
		handler.GetSong(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var song model.Song
		err = json.Unmarshal(w.Body.Bytes(), &song)
		assert.NoError(t, err)
		assert.Equal(t, patchedSongTitle, song.Title)
		assert.Equal(t, patchedSongTrackNumber, song.TrackNumber)
		assert.Equal(t, patchedDurationSeconds, song.DurationSeconds)
		assert.Equal(t, createdSong.ID, song.ID)

	})

	t.Run("DeleteSong", func(t *testing.T) {

		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var deleteSongUrl = "/songs/" + strconv.Itoa(createdSong.ID)
		c.Request = httptest.NewRequest("DELETE", deleteSongUrl, nil)
		c.Request.Header.Set("Content-Type", "application/json")

		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdSong.ID)}}

		// Call handler
		handler.DeleteSong(c)

		// Assertions
		assert.Equal(t, http.StatusNoContent, w.Code)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		// Verify the delete happened by doing a get request
		var getSongUrl = "/songs/" + strconv.Itoa(createdSong.ID)
		c.Request = httptest.NewRequest("GET", getSongUrl, nil)

		// Set the URL parameter that the handler expects
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdSong.ID)}}

		// Call handler
		handler.GetSong(c)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)

	})
}
