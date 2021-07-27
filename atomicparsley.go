package atomicparsley

import (
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const osType = runtime.GOOS

var (
	atomicPath string
	initErr    error
	atomicArgs = []string{
		"advisory",
		"album",
		"albumArtist",
		"apID",
		"artist",
		"artwork",
		"bpm",
		"category",
		"cnID",
		"comment",
		"compilation",
		"composer",
		"contentRating",
		"copyright",
		"description",
		"disk",
		"encodedBy",
		"encodingTool",
		"gapless",
		"geID",
		"genre",
		"grouping",
		"hdvideo",
		"keyword",
		"longdesc",
		"lyrics",
		"lyricsFile",
		"podcastGUID",
		"podcastURL",
		"productFlag",
		"purchaseDate",
		"stik",
		"storedesc",
		"title",
		"tracknum",
		"TVEpisode",
		"TVEpisodeNum",
		"TVNetwork",
		"TVSeasonNum",
		"TVShowName",
		"xID",
		"year",
	}
	config = map[string]map[string]string{
		"windows": {
			"atomicPath": filepath.Join(os.Getenv("TMP"), "AtomicParsley.exe"),
			"filename":   "AtomicParsleyWindows.exe",
		},
		"linux": {
			"atomicPath": filepath.Join("var", "tmp", "AtomicParsley"),
			"filename":   "AtomicParsleyLinux",
		},
		"darwin": {
			"atomicPath": filepath.Join("var", "tmp", "AtomicParsley"),
			"filename":   "AtomicParsleyMacOS",
		},
	}
)

func fileExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		return !f.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func downloadBinary(filename string) error {
	f, err := os.Create(atomicPath)
	if err != nil {
		return err
	}
	url := "https://github.com/Sorrow446/go-atomicparsley/releases/download/Bins/" + filename
	req, err := http.NewRequest(
		http.MethodGet, url, nil,
	)
	if err != nil {
		return err
	}
	req.Header.Set(
		"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "+
			"(KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
	)
	do, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return errors.New(do.Status)
	}
	_, err = io.Copy(f, do.Body)
	return err
}

func setup(cfg map[string]string) error {
	atomicPath = cfg["atomicPath"]
	exists, err := fileExists(atomicPath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	err = downloadBinary(cfg["filename"])
	if err != nil {
		return err
	}
	if osType != "windows" {
		err = os.Chmod(atomicPath, 0755)
	}
	return err
}

func filterTags(tags map[string]string) map[string]string {
	filteredTags := map[string]string{}
	for k, v := range tags {
		for _, arg := range atomicArgs {
			if arg == k {
				filteredTags[k] = v
				break
			}
		}
	}
	return filteredTags
}

func checkInput(path string, tags map[string]string) (map[string]string, error) {
	if len(tags) == 0 {
		return nil, errors.New("Tag map is empty.")
	}
	exists, err := fileExists(path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("Input file does not exist: " + path)
	}
	tags = filterTags(tags)
	if len(tags) == 0 {
		return nil, errors.New("All tags were filtered.")
	}
	coverValue, ok := tags["artwork"]
	if ok {
		exists, err := fileExists(coverValue)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("Input cover file does not exist: " + coverValue)
		}
	}
	return tags, err
}

func WriteTags(path string, tags map[string]string) error {
	if initErr != nil {
		return initErr
	}
	tags, err := checkInput(path, tags)
	if err != nil {
		return err
	}
	args := []string{path}
	base := atomicPath
// 	if osType != "windows" {
//  		base = "." + base
// 	}
	for k, v := range tags {
		args = append(args, "--"+k, v)
	}
	args = append(args, "-W")
	cmd := exec.Command(base, args...)
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	return err
}

func init() {
	cfg, ok := config[osType]
	if !ok {
		initErr = errors.New("Unsupported OS.")
	} else {
		initErr = setup(cfg)
	}
}
