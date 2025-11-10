# Setup and Run

## Start the docker
`docker compose up -d`

## Start the server
`DATABASE_URL=postgresql://postgres:gizzard@localhost:5432/albums go run main.go`


# Examples

## Example Artist Create
```
curl --request POST \
  --url http://localhost:9000/artists \
  --header 'Content-Type: application/json' \
  --data '{
  "name": "Nine Inch Nails",
  "description": "American industrial rock band formed by singer-songwriter, multi-instrumentalist, and producer Trent Reznor in 1988 in Cleveland, Ohio"
}'
```

## Example Album Create

```
curl --request POST \
  --url http://localhost:9000/albums \
  --header 'Content-Type: application/json' \
  --data '  {
    "artist_id": 5,
    "name": "The Fragile",
    "release_year": 1999
  },'
```

## Example Song Create

```
curl --request POST \
  --url http://localhost:9000/songs \
  --header 'Content-Type: application/json' \
  --data '  {
    "album_id": 9,
    "title": "Somewhat Damaged",
    "track_number": 1,
    "duration_seconds": 271
  }'
```

## Example Song Update

```
curl --request PUT \
  --url http://localhost:9000/songs/9 \
  --header 'Content-Type: application/json' \
  --data '  {
    "title": "Somewhat Damaged",
    "track_number": 2,
    "duration_seconds": 271
  }'
```

## Example Delete of Artist
Note: Deleting an artist will cascade delete the artist, album, and songs
```
curl --request DELETE \
  --url http://localhost:9000/artists/5 \
  --header 'Content-Type: application/json' \
  --data '  {
    "title": "Somewhat Damaged",
    "track_number": 1,
    "duration_seconds": 271
  }'
```


