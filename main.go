package main

import (
	"context"
	"errors"
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

type CreateArtist struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateArtist struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type PatchArtist struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type Album struct {
	ID          int    `json:"id"`
	ArtistName  string `json:"artist,omitempty"`
	Name        string `json:"name"`
	ReleaseYear int    `json:"release_year"`
}

type CreateAlbum struct {
	ArtistID    int    `json:"artist_id"`
	Name        string `json:"name"`
	ReleaseYear int    `json:"release_year"`
}

type UpdateAlbum struct {
	Name        string `json:"name"`
	ReleaseYear int    `json:"release_year"`
}

type PatchAlbum struct {
	Name        *string `json:"name"`
	ReleaseYear *int    `json:"release_year"`
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

type SongResponse struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	TrackNumber     int    `json:"track_number"`
	DurationSeconds int    `json:"duration_seconds"`
}

type CreateSong struct {
	AlbumID         int    `json:"album_id"`
	Title           string `json:"title"`
	TrackNumber     int    `json:"track_number"`
	DurationSeconds int    `json:"duration_seconds"`
}

type UpdateSong struct {
	Title           string `json:"title"`
	TrackNumber     int    `json:"track_number"`
	DurationSeconds int    `json:"duration_seconds"`
}

type PatchSong struct {
	Title           *string `json:"title"`
	TrackNumber     *int    `json:"track_number"`
	DurationSeconds *int    `json:"duration_seconds"`
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
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return artists, nil
}

func getArtist(c *gin.Context, id string) (*ArtistWithAlbums, error) {
	var artist ArtistWithAlbums
	query := `SELECT id, name, description FROM artist WHERE id = $1`

	err := dbPool.QueryRow(c, query, id).Scan(&artist.ID, &artist.Name, &artist.Description)
	if err != nil {
		return nil, err
	}

	artist.Albums, err = getAlbumsForArtist(c, artist.ID)
	if err != nil {
		return nil, err
	}

	return &artist, nil

}

func createArtist(c *gin.Context, artist CreateArtist) error {

	query := `INSERT INTO artist (name, description) VALUES ($1, $2)`
	_, err := dbPool.Query(c, query, artist.Name, artist.Description)
	if err != nil {
		return err
	}

	return nil

}

func updateArtist(c *gin.Context, artist UpdateArtist, id string) error {
	query := `UPDATE artist SET name = $2, description = $3 WHERE id = $1`

	_, err := dbPool.Query(c, query, id, artist.Name, artist.Description)
	if err != nil {
		return err
	}

	return nil

}

func patchArtist(c *gin.Context, artist PatchArtist, id string) error {

	if artist.Name != nil {
		queryName := `UPDATE artist SET name = $2 WHERE id = $1`
		_, err := dbPool.Query(c, queryName, id, artist.Name)
		if err != nil {
			return err
		}
	}

	if artist.Description != nil {
		queryDescription := `UPDATE artist SET description = $2 WHERE id = $1`
		_, err := dbPool.Query(c, queryDescription, id, artist.Description)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteArtist(c *gin.Context, id string) error {

	query := `DELETE FROM artist where id = $1`

	_, err := dbPool.Query(c, query, id)
	if err != nil {
		return err
	}

	return nil
}

func getAlbums(c *gin.Context) ([]Album, error) {

	query := `SELECT album.id, album.name, album.release_year, artist.name as artist 
			FROM album 
			JOIN artist ON album.artist_id = artist.id 
			ORDER BY artist.name, album.name`
	rows, err := dbPool.Query(c, query)
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
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

func getAlbum(c *gin.Context, id string) (*AlbumWithSongs, error) {

	var album AlbumWithSongs

	query := `SELECT album.id, album.name, album.release_year, artist.name as artist 
		FROM album 
		JOIN artist ON album.artist_id = artist.id 
		WHERE album.id = $1`

	err := dbPool.QueryRow(c, query, id).Scan(&album.ID, &album.Name, &album.ReleaseYear, &album.ArtistName)
	if err != nil {
		return nil, err
	}

	album.Songs, err = getSongsForAlbum(c, album.ID)
	if err != nil {
		return nil, err
	}

	return &album, nil

}

func createAlbum(c *gin.Context, album CreateAlbum) error {

	query := `INSERT INTO album (artist_id, name, release_year) VALUES ($1, $2, $3)`
	_, err := dbPool.Query(c, query, album.ArtistID, album.Name, album.ReleaseYear)
	if err != nil {
		return err
	}

	return nil

}

func updateAlbum(c *gin.Context, album UpdateAlbum, id string) error {
	query := `UPDATE album SET name = $2, release_year = $3 WHERE id = $1`

	_, err := dbPool.Query(c, query, id, album.Name, album.ReleaseYear)
	if err != nil {
		return err
	}

	return nil
}

func patchAlbum(c *gin.Context, album PatchAlbum, id string) error {

	if album.Name != nil {
		queryName := `UPDATE album SET name = $2 WHERE id = $1`
		_, err := dbPool.Query(c, queryName, id, album.Name)
		if err != nil {
			return err
		}
	}

	if album.ReleaseYear != nil {
		queryReleaseYear := `UPDATE album SET release_year = $2 WHERE id = $1`
		_, err := dbPool.Query(c, queryReleaseYear, id, album.ReleaseYear)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteAlbum(c *gin.Context, id string) error {

	query := `DELETE FROM album where id = $1`

	_, err := dbPool.Query(c, query, id)
	if err != nil {
		return err
	}

	return nil
}

func getAlbumsForArtist(c *gin.Context, artistID int) ([]Album, error) {

	query := `SELECT id, name, release_year FROM album WHERE artist_id = $1 ORDER BY release_year DESC`
	rows, err := dbPool.Query(c, query, artistID)
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
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

func getSongs(c *gin.Context) ([]Song, error) {

	query := `SELECT song.id, song.title, song.track_number, song.duration_seconds, album.name as album, artist.name as artist
				FROM song
				JOIN album ON song.album_id = album.id
				JOIN artist ON album.artist_id = artist.id
				ORDER BY song.title`
	rows, err := dbPool.Query(c, query)
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
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func getSong(c *gin.Context, id string) (*Song, error) {

	query := `SELECT song.id, song.title, song.track_number, song.duration_seconds, album.name as album, artist.name as artist
				FROM song
				JOIN album ON song.album_id = album.id
				JOIN artist ON album.artist_id = artist.id
				WHERE song.id = $1`
	var song Song

	err := dbPool.QueryRow(c, query, id).Scan(&song.ID, &song.Title, &song.TrackNumber, &song.DurationSeconds, &song.AlbumName, &song.ArtistName)
	if err != nil {
		return nil, err
	}

	return &song, nil
}

func createSong(c *gin.Context, song CreateSong) (*SongResponse, error) {

	var songCreated SongResponse

	query := `INSERT INTO song (album_id, title, track_number, duration_seconds) VALUES ($1, $2, $3, $4) RETURNING id, title, track_number, duration_seconds`
	err := dbPool.QueryRow(c, query, song.AlbumID, song.Title, song.TrackNumber, song.DurationSeconds).Scan(&songCreated.ID, &songCreated.Title, &songCreated.TrackNumber, &songCreated.DurationSeconds)
	if err != nil {
		return nil, err
	}

	return &songCreated, nil

}

func updateSong(c *gin.Context, song UpdateSong, id string) error {
	query := `UPDATE song SET title = $2, track_number = $3, duration_seconds = $4 WHERE id = $1`

	_, err := dbPool.Query(c, query, id, song.Title, song.TrackNumber, song.DurationSeconds)
	if err != nil {
		return err
	}

	return nil
}

func patchSong(c *gin.Context, song PatchSong, id string) error {

	if song.Title != nil {
		queryName := `UPDATE song SET title = $2 WHERE id = $1`
		_, err := dbPool.Query(c, queryName, id, song.Title)
		if err != nil {
			return err
		}
	}

	if song.TrackNumber != nil {
		queryReleaseYear := `UPDATE song SET track_number = $2 WHERE id = $1`
		_, err := dbPool.Query(c, queryReleaseYear, id, song.TrackNumber)
		if err != nil {
			return err
		}
	}

	if song.DurationSeconds != nil {
		queryReleaseYear := `UPDATE song SET duration_seconds = $2 WHERE id = $1`
		_, err := dbPool.Query(c, queryReleaseYear, id, song.DurationSeconds)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteSong(c *gin.Context, id string) error {

	query := `DELETE FROM song where id = $1`

	_, err := dbPool.Query(c, query, id)
	if err != nil {
		return err
	}

	return nil
}

func getSongsForAlbum(c *gin.Context, albumID int) ([]Song, error) {

	query := `SELECT id, title, track_number, duration_seconds FROM song WHERE album_id = $1 ORDER BY track_number`
	rows, err := dbPool.Query(c, query, albumID)
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
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
			log.Printf("Error fetching artists: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, artists)
	})

	router.GET("/artists/:id", func(c *gin.Context) {
		id := c.Param("id")
		artist, err := getArtist(c, id)

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
	})

	router.DELETE("/artists/:id", func(c *gin.Context) {
		id := c.Param("id")
		err := deleteArtist(c, id)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
				return
			}

			log.Printf("Error deleting artist %s: %v", id, err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, "Deleted Artist")
	})

	router.POST("/artists", func(c *gin.Context) {
		var newArtist CreateArtist

		err := c.BindJSON(&newArtist)
		if err != nil {
			return
		}

		err = createArtist(c, newArtist)

		if err != nil {
			log.Printf("Error creating artist %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusCreated, "Created Artist")
	})

	router.PUT("/artists/:id", func(c *gin.Context) {
		id := c.Param("id")
		var newArtist UpdateArtist

		err := c.BindJSON(&newArtist)
		if err != nil {
			return
		}

		err = updateArtist(c, newArtist, id)

		if err != nil {
			log.Printf("Error updating artist %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, "Updated Artist")
	})

	router.PATCH("/artists/:id", func(c *gin.Context) {
		id := c.Param("id")

		var newArtist PatchArtist

		err := c.BindJSON(&newArtist)
		if err != nil {
			return
		}

		err = patchArtist(c, newArtist, id)

		if err != nil {
			log.Printf("Error patching artist %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, "Patched Artist")
	})

	router.GET("/albums", func(c *gin.Context) {
		albums, err := getAlbums(c)

		if err != nil {
			log.Printf("Error fetching artists: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, albums)
	})

	router.GET("/albums/:id", func(c *gin.Context) {
		id := c.Param("id")
		album, err := getAlbum(c, id)

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
	})

	router.POST("/albums", func(c *gin.Context) {
		var newAlbum CreateAlbum

		err := c.BindJSON(&newAlbum)
		if err != nil {
			return
		}

		err = createAlbum(c, newAlbum)

		if err != nil {
			log.Printf("Error creating album %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusCreated, "Created Album")
	})

	router.PUT("/albums/:id", func(c *gin.Context) {
		id := c.Param("id")
		var newAlbum UpdateAlbum

		err := c.BindJSON(&newAlbum)
		if err != nil {
			return
		}

		err = updateAlbum(c, newAlbum, id)

		if err != nil {
			log.Printf("Error updating album %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, "Updated Album")
	})

	router.PATCH("/albums/:id", func(c *gin.Context) {
		id := c.Param("id")
		var newAlbum PatchAlbum

		err := c.BindJSON(&newAlbum)
		if err != nil {
			return
		}

		err = patchAlbum(c, newAlbum, id)

		if err != nil {
			log.Printf("Error patching album %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, "Patched Album")
	})

	router.DELETE("/albums/:id", func(c *gin.Context) {
		id := c.Param("id")
		err := deleteAlbum(c, id)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
				return
			}

			log.Printf("Error deleting album %s: %v", id, err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, "Deleted Album")
	})

	router.GET("/songs", func(c *gin.Context) {
		songs, err := getSongs(c)

		if err != nil {
			log.Printf("Error fetching songs: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, songs)
	})

	router.POST("/songs", func(c *gin.Context) {
		var newSong CreateSong

		err := c.BindJSON(&newSong)
		if err != nil {
			return
		}

		songCreated, err := createSong(c, newSong)

		if err != nil {
			log.Printf("Error creating song %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusCreated, songCreated)
	})

	router.GET("/songs/:id", func(c *gin.Context) {
		id := c.Param("id")
		song, err := getSong(c, id)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
				return
			}

			log.Printf("Error fetching song %s: %v", id, err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, song)
	})

	router.PUT("/songs/:id", func(c *gin.Context) {
		id := c.Param("id")
		var newSong UpdateSong

		err := c.BindJSON(&newSong)
		if err != nil {
			return
		}

		err = updateSong(c, newSong, id)

		if err != nil {
			log.Printf("Error updating song %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, "Updated Song")
	})

	router.PATCH("/songs/:id", func(c *gin.Context) {
		id := c.Param("id")
		var newSong PatchSong

		err := c.BindJSON(&newSong)
		if err != nil {
			return
		}

		err = patchSong(c, newSong, id)

		if err != nil {
			log.Printf("Error patching song %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, "Patched Song")
	})

	router.DELETE("/songs/:id", func(c *gin.Context) {
		id := c.Param("id")
		err := deleteSong(c, id)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
				return
			}

			log.Printf("Error deleting song %s: %v", id, err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, "Deleted Song")
	})

	router.Run(":9000")
}
