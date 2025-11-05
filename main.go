package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// struct to represent artist
type artist struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// var artists = []artist{
// 	{ID: "1", Name: "King Gizzard & the Lizard Wizard", Description: "Lizard Wizard or something idk."},
// 	{ID: "2", Name: "Minus The Bear", Description: "Band that doesn't have a bear."},
// 	{ID: "3", Name: "Kraftwerk", Description: "Idk"},
// }

// func getArtists2(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, artists)
// }

func getArtists(c *gin.Context) ([]artist, error) {
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

	artists := []artist{}

	for rows.Next() {
		var art artist
		if err := rows.Scan(&art.ID, &art.Name, &art.Description); err != nil {
			return nil, err
		}

		artists = append(artists, art)

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}

	return artists, nil
}

func main() {
	// rows, _ := conn.Query(context.Background(), "SELECT id, name, description FROM artist", 5)
	// if err != nil {
	// 	fmt.Println("ERROR:", err)
	// }

	// fmt.Println(artists)

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

	router.Run(":9000")
}
