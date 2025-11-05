package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

type Artist struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Album struct {
	ID          int    `json:"id"`
	ArtistName  string `json:"artist,omitempty"`
	Name        string `json:"name"`
	ReleaseYear int    `json:"release_year"`
}

type ArtistWithAlbums struct {
	Artist
	Albums []Album
}

type Song struct {
	ID              int    `json:"id"`
	ArtistName      string `json:"artist,omitempty"`
	AlbumName       string `json:"album,omitempty"`
	Title           string `json:"title"`
	TrackNumber     int    `json:"track_number"`
	DurationSeconds int    `json:"duration_seconds"`
}

type AlbumWithSongs struct {
	Album
	Songs []Song
}

func getArtists(c *gin.Context) ([]Artist, error) {
	query := `SELECT id, name, description FROM artist`
	rows, err := dbPool.Query(c, query)
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

	query := `SELECT album.id, album.name, album.release_year, artist.name as artist 
			FROM album 
			JOIN artist ON album.artist_id = artist.id 
			ORDER BY artist.name, album.name`
	rows, err := conn.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	albums := []Album{}

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Name, &album.ReleaseYear, &album.ArtistName); err != nil {
			return nil, err
		}

		albums = append(albums, album)

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}

	return albums, nil
}

func getAlbum(c *gin.Context, id string) (*AlbumWithSongs, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	var album AlbumWithSongs

	query := `SELECT album.id, album.name, album.release_year, artist.name as artist 
		FROM album 
		JOIN artist ON album.artist_id = artist.id 
		WHERE album.id = $1`

	err = conn.QueryRow(c, query, id).Scan(&album.ID, &album.Name, &album.ReleaseYear, &album.ArtistName)
	if err != nil {
		return nil, err
	}

	album.Songs, err = getSongsForAlbum(c, album.ID)
	if err != nil {
		return nil, err
	}

	return &album, err

}

func getAlbumsForArtist(c *gin.Context, artistID int) ([]Album, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `SELECT id, name, release_year FROM album WHERE artist_id = $1 ORDER BY release_year DESC`
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

func getSongs(c *gin.Context) ([]Song, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `SELECT song.id, song.title, song.track_number, song.duration_seconds, album.name as album, artist.name as artist
				FROM song
				JOIN album ON song.album_id = album.id
				JOIN artist ON album.artist_id = artist.id
				ORDER BY song.title`
	rows, err := conn.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := []Song{}

	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.ID, &song.Title, &song.TrackNumber, &song.DurationSeconds, &song.AlbumName, &song.ArtistName); err != nil {
			return nil, err
		}

		songs = append(songs, song)

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}

	return songs, nil
}

func getSongsForAlbum(c *gin.Context, albumID int) ([]Song, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `SELECT id, title, track_number, duration_seconds FROM song WHERE album_id = $1 ORDER BY track_number`
	rows, err := conn.Query(c, query, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := []Song{}

	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.ID, &song.Title, &song.TrackNumber, &song.DurationSeconds); err != nil {
			return nil, err
		}

		songs = append(songs, song)

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}

	return songs, nil
}

func main() {

	var err error
	dbPool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}
	defer dbPool.Close()

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

	router.GET("/songs", func(c *gin.Context) {
		songs, err := getSongs(c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, songs)
	})

	router.Run(":9000")
}
