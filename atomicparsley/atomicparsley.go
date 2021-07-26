package atomicparsley

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const osType = runtime.GOOS

var atomicPath string

func fileExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		return !f.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func download(path, fname string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	url := "https://github.com/wez/atomicparsley/releases/download/20210715.151551.e7ad03a/" + fname
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
	bodyBytes, err := ioutil.ReadAll(do.Body)
	if err != nil {
		return err
	}
	zipReader, err := zip.NewReader(bytes.NewReader(bodyBytes), int64(len(bodyBytes)))
	if err != nil {
		log.Fatal(err)
	}
	zf, err := zipReader.File[0].Open()
	if err != nil {
		return nil
	}
	defer zf.Close()
	_, err = io.Copy(f, zf)
	return err
}

func winSetup() error {
	atomicPath = filepath.Join(os.Getenv("TMP"), "AtomicParsley.exe")
	exists, err := fileExists(atomicPath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	fname := "AtomicParsleyWindows.zip"
	err = download(atomicPath, fname)
	return err
}

func linuxSetup() error {
	var fname string
	atomicPath = filepath.Join("/var/tmp/", "AtomicParsley")
	exists, err := fileExists(atomicPath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if osType == "linux" {
		fname = "AtomicParsleyLinux.zip"
	} else {
		fname = "AtomicParsleyMacOS.zip"
	}
	err = download(atomicPath, fname)
	if err != nil {
		return err
	}
	err = os.Chmod(atomicPath, 0755)
	return err
}

func init() {
	switch os := osType; os {
	case "windows":
		winSetup()
	case "linux":
		linuxSetup()
	default:
		panic("Unsupported OS.")
	}
}

func writeTags(path string, tags map[string]string) error {
	args := []string{path}
	base := atomicPath
	if osType == "linux" {
		base = "./" + base
	}
	for k, v := range tags {
		args = append(args, "--"+k, v)
	}
	args = append(args, "-W")
	cmd := exec.Command(base, args...)
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	return err
}
