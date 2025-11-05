package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Artist struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Album struct {
	ID int `json:"id"`
	// ArtistID    int    `json:"artist_id"`
	Name        string `json:"name"`
	ReleaseYear int    `json:"release_year"`
}

type ArtistWithAlbums struct {
	Artist
	Albums []Album
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

func getArtist(c *gin.Context, id string) (*ArtistWithAlbums, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	var artist ArtistWithAlbums
	query := `SELECT id, name, description FROM artist WHERE id = $1`

	err = conn.QueryRow(c, query, id).Scan(&artist.ID, &artist.Name, &artist.Description)
	if err != nil {
		return nil, err
	}

	artist.Albums, err = getAlbumsForArtist(c, artist.ID)
	if err != nil {
		return nil, err
	}

	return &artist, err

}

func getAlbums(c *gin.Context) ([]Album, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `SELECT id, name, release_year FROM album`
	rows, err := conn.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	albums := []Album{}

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Name, &album.ReleaseYear); err != nil {
			return nil, err
		}

		albums = append(albums, album)

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}

	return albums, nil
}

func getAlbum(c *gin.Context, id string) (*Album, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	var album Album
	query := `SELECT id, name, release_year FROM album WHERE id = $1`

	err = conn.QueryRow(c, query, id).Scan(&album.ID, &album.Name, &album.ReleaseYear)
	if err != nil {
		return nil, err
	}

	// artist.Albums, err = getAlbumsForArtist(c, artist.ID)
	// if err != nil {
	// 	return nil, err
	// }

	return &album, err

}

func getAlbumsForArtist(c *gin.Context, artistID int) ([]Album, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `SELECT id, name, release_year FROM album WHERE artist_id = $1`
	rows, err := conn.Query(c, query, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	albums := []Album{}

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Name, &album.ReleaseYear); err != nil {
			return nil, err
		}

		albums = append(albums, album)

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}

	return albums, nil
}

func main() {

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {

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

		c.JSON(http.StatusOK, artists)
	})

	router.GET("/artists/:id", func(c *gin.Context) {
		id := c.Param("id")
		artist, err := getArtist(c, id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, artist)
	})

	router.GET("/albums", func(c *gin.Context) {
		albums, err := getAlbums(c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, albums)
	})

	router.GET("/albums/:id", func(c *gin.Context) {
		id := c.Param("id")
		album, err := getAlbum(c, id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, album)
	})

	router.Run(":9000")
}
