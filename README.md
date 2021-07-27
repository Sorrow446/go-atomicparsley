# go-atomicparsley
Go wrapper for the mp4/m4a tagger AtomicParsley.

# Setup
```go
import ap "github.com/Sorrow446/go-atomicparsley"
```
The appropriate binary for your OS will be automatically fetched on the first start-up (250-500KB).

# Usage
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
