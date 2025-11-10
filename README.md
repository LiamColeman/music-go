# Music API
REST API for artists, albums, and songs backed by a PostgreSQL database. This is a sample project I build for learning go.

## Running
Start docker, then start the server

    docker compose up -d
    DATABASE_URL=postgresql://postgres:gizzard@localhost:5432/albums go run main.go

## Enable Live Reload During Development (Optional)

Install [Air](https://github.com/air-verse/air?tab=readme-ov-file#via-go-install-recommended)

    go install github.com/air-verse/air@latest

Start Air

    Air



## Examples

Create an artist
```
curl --request POST \
  --url http://localhost:8080/artists \
  --header 'Content-Type: application/json' \
  --data '{
  "name": "Nine Inch Nails",
  "description": "American industrial rock band formed by singer-songwriter, multi-instrumentalist, and producer Trent Reznor in 1988 in Cleveland, Ohio"
}'
```

Create an album
```
curl --request POST \
  --url http://localhost:8080/albums \
  --header 'Content-Type: application/json' \
  --data '  {
    "artist_id": 5,
    "name": "The Fragile",
    "release_year": 1999
  },'
```

Create a song
```
curl --request POST \
  --url http://localhost:8080/songs \
  --header 'Content-Type: application/json' \
  --data '  {
    "album_id": 9,
    "title": "Somewhat Damaged",
    "track_number": 1,
    "duration_seconds": 271
  }'
```

Update a song
```
curl --request PUT \
  --url http://localhost:8080/songs/9 \
  --header 'Content-Type: application/json' \
  --data '  {
    "title": "Somewhat Damaged",
    "track_number": 2,
    "duration_seconds": 271
  }'
```

Delete an Artist 
Note: Deleting an artist will cascade delete the albums and songs for the artist
```
curl --request DELETE \
  --url http://localhost:8080/artists/5 \
  --header 'Content-Type: application/json' \
  --data '  {
    "title": "Somewhat Damaged",
    "track_number": 1,
    "duration_seconds": 271
  }'
```


