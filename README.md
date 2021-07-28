# go-atomicparsley
Go wrapper for the mp4/m4a tagger AtomicParsley for Windows, Linux and macOS.

# Setup
```go
import ap "github.com/Sorrow446/go-atomicparsley"
```
The appropriate binary for your OS will be automatically fetched on the first start-up (250-500KB).

# Usage

## Read
```go
ReadTags(path string) (map[string]string, error)
```
```go
	tags, err := ap.ReadTags("17.m4a")
	if err != nil {
		panic(err)
	}
	for k, v := range tags {
		fmt.Println(k+":", v)
	}
```
Read all tags and print them.

Output:
```
ISRC: FISGC9900001
copyright: ℗ 2000 16 Inch Records
UPC: 743217881122
album: Before the Storm
year: 2001
genre: Dance
compilation: false
LABEL: 16 Inch Records
tracknum: 1 of 11
composer: Jaakko Salovaara & Ville Virtanen
```

## Write
```go
WriteTags(path string, tags map[string]string) error
```

```go
	tags := map[string]string{
		"album":       "Before the Storm",
		"albumArtist": "Darude",
		"artist":      "Darude",
		"artwork":     "cover.jpg",
		"track":       "Feel The Beat",
		"tracknum":    "3/15",
		"year":        "2000",
	}
	err := ap.WriteTags("in.m4a", tags)
	if err != nil {
		panic(err)
	}
```
Write album, album artist, artist, track, tracknum/total, year tags and add cover from "cover.jpg".

```go
	tags := map[string]string{
		"album":  "",
	}
	err := ap.WriteTags("in.m4a", tags)
	if err != nil {
		panic(err)
	}
```
Delete album tag.

```go
	tags := map[string]string{
		"artwork":     "REMOVE_ALL"
	}
	err := ap.WriteTags("in.m4a", tags)
	if err != nil {
		panic(err)
	}
```
Remove all covers.


# Keys
## Read
```
album
albumArtist
artist
artistSort
bpm
comment
compilation
composer
composerSort
conductor
contentGroup
copyright
description
director
disk
genre
genre
itunesAccount
itunesAdvisory
itunesAlbumId
itunesArtistId
itunesCatalogId
itunesComposerId
itunesCountryId
itunesGapless
itunesGenreId
itunesHdVideo
itunesMediaType
itunesOwner
itunesPurchaseDate
lyrics
movement
movementName
movementTotal
podcast
podcastCategory
podcastDesc
podcastId
podcastKeywords
podcastUrl
title
titleSort
tracknum
tvEpisode
tvEpisodeId
tvNetwork
tvSeason
tvShow
tvShowSort
work
xID
year
```
Customs tags are also supported. The keys for them will be in uppercase.

## Write
```
advisory
album
albumArtist
apID
artist
artwork
bpm
category
cnID
comment
compilation
composer
contentRating
copyright
description
disk
encodedBy
encodingTool
gapless
geID
genre
grouping
hdvideo
keyword
longdesc
lyrics
lyricsFile
podcastGUID
podcastURL
productFlag
purchaseDate
stik
storedesc
title
tracknum
TVEpisode
TVEpisodeNum
TVNetwork
TVSeasonNum
TVShowName
xID
year
```
