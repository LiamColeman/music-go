package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// struct to represent artist
type Artist struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func getArtists(c *gin.Context) ([]Artist, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `SELECT id, name, description FROM artist`
	rows, err := conn.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	artists := []Artist{}

	for rows.Next() {
		var artist Artist
		if err := rows.Scan(&artist.ID, &artist.Name, &artist.Description); err != nil {
			return nil, err
		}

		artists = append(artists, artist)

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}

	return artists, nil
}

func getArtist(c *gin.Context, id string) (*Artist, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `SELECT id, name, description FROM artist WHERE id = $1`
	rows, err := conn.Query(c, query, id)
	if err != nil {
		return nil, err
	}

	artist, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Artist])
	if err != nil {
		return nil, err
	}

	return &artist, err

}

func main() {

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		// Return JSON Response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/artists", func(c *gin.Context) {
		artists, err := getArtists(c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return JSON Response
		c.JSON(http.StatusOK, artists)
	})

	router.GET("/artists/:id", func(c *gin.Context) {
		id := c.Param("id")
		artist, err := getArtist(c, id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return JSON Response
		c.JSON(http.StatusOK, artist)
	})

	router.Run(":9000")
}
