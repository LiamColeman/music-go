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

// setupTestDBArtist creates a connection pool to your test database
// NOTE: Start test database with: docker compose -f docker-compose.test.yml up -d
func setupTestDBArtist(t *testing.T) *pgxpool.Pool {
	// Use test database URL (separate DB on port 5433)
	databaseURL := "postgresql://postgres:gizzard@localhost:5433/albums_test"

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Fatalf("Unable to create connection pool: %v", err)
	}

	return pool
}

// cleanupArtists removes test data from the database
func cleanupArtists(t *testing.T, pool *pgxpool.Pool, name string) {
	_, err := pool.Exec(context.Background(),
		"DELETE FROM artist WHERE name = $1", name)
	if err != nil {
		t.Logf("Warning: cleanup failed: %v", err)
	}
}

func TestArtists(t *testing.T) {
	// Setup
	pool := setupTestDBArtist(t)
	defer pool.Close()

	artistRepo := repository.NewArtistRepository(pool)
	handler := NewArtistHandler(artistRepo)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	testArtistName := "Test Artist for Integration"
	testArtistDescription := "A test artist for integration testing"
	defer cleanupArtists(t, pool, testArtistName)

	// Store this so we can use in tests
	var createdArtist model.Artist

	t.Run("CreateArtist", func(t *testing.T) {
		// Create request body
		createArtist := model.CreateArtist{
			Name:        testArtistName,
			Description: testArtistDescription,
		}

		body, _ := json.Marshal(createArtist)

		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/artists", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call handler
		handler.CreateArtist(c)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &createdArtist)
		assert.NoError(t, err)
		assert.Equal(t, testArtistName, createdArtist.Name)
		assert.Equal(t, testArtistDescription, createdArtist.Description)
		assert.NotZero(t, createdArtist.ID)

		// Check Location header
		location := w.Header().Get("Location")
		locationUrl := "/artists/" + strconv.Itoa(createdArtist.ID)
		assert.Contains(t, location, locationUrl)
	})

	t.Run("GetArtist", func(t *testing.T) {
		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var getArtistUrl = "/artists/" + strconv.Itoa(createdArtist.ID)
		c.Request = httptest.NewRequest("GET", getArtistUrl, nil)

		// Set the URL parameter that the handler expects
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdArtist.ID)}}

		// Call handler
		handler.GetArtist(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var artist model.Artist
		err := json.Unmarshal(w.Body.Bytes(), &artist)
		assert.NoError(t, err)
		assert.Equal(t, testArtistName, artist.Name)
		assert.Equal(t, testArtistDescription, artist.Description)
		assert.Equal(t, createdArtist.ID, artist.ID)
	})

	t.Run("GetAllArtists", func(t *testing.T) {
		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/artists", nil)

		// Call handler
		handler.GetAll(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var artists []model.Artist
		err := json.Unmarshal(w.Body.Bytes(), &artists)
		assert.NoError(t, err)
		assert.NotEmpty(t, artists)

		// Verify our test artist is in the list
		found := false
		for _, artist := range artists {
			if artist.ID == createdArtist.ID {
				found = true
				assert.Equal(t, testArtistName, artist.Name)
				assert.Equal(t, testArtistDescription, artist.Description)
				break
			}
		}
		assert.True(t, found, testArtistDescription)
	})

	t.Run("UpdateArtist", func(t *testing.T) {

		updatedArtistName := "Update Artist for Integration"
		updatedArtistDescription := "An updated test artist for integration testing"

		// Store this so we can use in tests
		var updatedArtist model.Artist

		// Create request body
		updateArtist := model.UpdateArtist{
			Name:        updatedArtistName,
			Description: updatedArtistDescription,
		}

		body, _ := json.Marshal(updateArtist)

		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var updateArtistUrl = "/artists/" + strconv.Itoa(createdArtist.ID)
		c.Request = httptest.NewRequest("PUT", updateArtistUrl, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdArtist.ID)}}

		// Call handler
		handler.UpdateArtist(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &updatedArtist)
		assert.NoError(t, err)
		assert.Equal(t, updatedArtistName, updatedArtist.Name)
		assert.Equal(t, updatedArtistDescription, updatedArtist.Description)
		assert.NotZero(t, updatedArtist.ID)

		// Check Location header
		location := w.Header().Get("Location")
		locationUrl := "/artists/" + strconv.Itoa(updatedArtist.ID)
		assert.Contains(t, location, locationUrl)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		// Verify the update happened by doing a get request
		var getArtistUrl = "/artists/" + strconv.Itoa(updatedArtist.ID)
		c.Request = httptest.NewRequest("GET", getArtistUrl, nil)

		// Set the URL parameter that the handler expects
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(updatedArtist.ID)}}

		// Call handler
		handler.GetArtist(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var artist model.Artist
		err = json.Unmarshal(w.Body.Bytes(), &artist)
		assert.NoError(t, err)
		assert.Equal(t, updatedArtistName, artist.Name)
		assert.Equal(t, updatedArtistDescription, artist.Description)
		assert.Equal(t, updatedArtist.ID, artist.ID)

	})

	t.Run("PatchArtist", func(t *testing.T) {

		patchedArtistName := "Patched Artist for Integration"
		patchedArtistDescription := "A patched test artist for integration testing"

		// Store this so we can use in tests
		var patchedArtist model.Artist

		// Create request body
		patchArtist := model.PatchArtist{
			Name:        &patchedArtistName,
			Description: &patchedArtistDescription,
		}

		body, _ := json.Marshal(patchArtist)

		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var updateArtistUrl = "/artists/" + strconv.Itoa(createdArtist.ID)
		c.Request = httptest.NewRequest("PATCH", updateArtistUrl, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdArtist.ID)}}

		// Call handler
		handler.PatchArtist(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &patchedArtist)
		assert.NoError(t, err)
		assert.Equal(t, patchedArtistName, patchedArtist.Name)
		assert.Equal(t, patchedArtistDescription, patchedArtist.Description)
		assert.NotZero(t, patchedArtist.ID)

		// Check Location header
		location := w.Header().Get("Location")
		locationUrl := "/artists/" + strconv.Itoa(patchedArtist.ID)
		assert.Contains(t, location, locationUrl)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		// Verify the update happened by doing a get request
		var getArtistUrl = "/artists/" + strconv.Itoa(patchedArtist.ID)
		c.Request = httptest.NewRequest("GET", getArtistUrl, nil)

		// Set the URL parameter that the handler expects
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(patchedArtist.ID)}}

		// Call handler
		handler.GetArtist(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var artist model.Artist
		err = json.Unmarshal(w.Body.Bytes(), &artist)
		assert.NoError(t, err)
		assert.Equal(t, patchedArtistName, artist.Name)
		assert.Equal(t, patchedArtistDescription, artist.Description)
		assert.Equal(t, patchedArtist.ID, artist.ID)

	})

	t.Run("DeleteArtist", func(t *testing.T) {

		// Create HTTP request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var deleteArtistUrl = "/artists/" + strconv.Itoa(createdArtist.ID)
		c.Request = httptest.NewRequest("DELETE", deleteArtistUrl, nil)
		c.Request.Header.Set("Content-Type", "application/json")

		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdArtist.ID)}}

		// Call handler
		handler.DeleteArtist(c)

		// Assertions
		assert.Equal(t, http.StatusNoContent, w.Code)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		// Verify the delete happened by doing a get request
		var getArtistUrl = "/artists/" + strconv.Itoa(createdArtist.ID)
		c.Request = httptest.NewRequest("GET", getArtistUrl, nil)

		// Set the URL parameter that the handler expects
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(createdArtist.ID)}}

		// Call handler
		handler.GetArtist(c)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)

	})
}
