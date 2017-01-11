package utils

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"path/filepath"

	"github.com/bulletind/moire/signature"
	"github.com/bulletind/spyder/log"
)

type mediaSettings struct {
	Host        string
	PublicKey   string
	SecretKey   string
	DownloadDir string
}

// load once, store for reuse
var settings *mediaSettings

func loadConfig() {
	if settings != nil {
		return
	}

	// 1st create dir
	dir := filepath.Join(os.TempDir(), "Downloads")
	os.MkdirAll(dir, 0777)

	settings = &mediaSettings{
		Host:        GetEnv("MOIRE_HOST", false),
		PublicKey:   GetEnv("MOIRE_PUBLIC_KEY", false),
		SecretKey:   GetEnv("MOIRE_SECRET_KEY", false),
		DownloadDir: dir,
	}
}

func DownloadFile(url string, extension string, isPrivate bool) (fileName string, size int64, err error) {
	loadConfig()

	fileName, err = getFilePath(url, extension)
	if err != nil {
		return
	}

	// simply check if file exists
	if _, err = os.Stat(fileName); err == nil {
		return
	}

	// so we need to download
	url = makeUrl(url, isPrivate)
	size, err = download(url, fileName)
	return
}

func makeUrl(url string, isPrivate bool) string {
	if isPrivate {
		if settings.PublicKey == "" {
			log.Error("When using private urls, you need to provide keys for the mediaserver")
		} else {
			// add host info when needed
			if !strings.HasPrefix(url, "http") {
				if !strings.HasSuffix(settings.Host, "/") && !strings.HasPrefix(url, "/") {
					url = "/" + url
				}
				url = settings.Host + url
			}
			url = signature.MakeUrl(settings.PublicKey, settings.SecretKey, url)
		}
	}
	return url
}

func getFilePath(rawURL string, extension string) (fileName string, err error) {
	fileURL, err := url.Parse(rawURL)

	if err != nil {
		return
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	// add date and hour path so we can re-use downloaded files
	fileName = filepath.Join(settings.DownloadDir, time.Now().Format("2006010215"))
	fileName = filepath.Join(fileName, fileURL.Host)
	for i := 1; i < len(segments)-1; i++ {
		fileName = filepath.Join(fileName, segments[i])
	}
	err = os.MkdirAll(fileName, 0777)
	if err != nil {
		return
	}

	name := segments[len(segments)-1]
	if !strings.Contains(name, ".") {
		if !strings.HasPrefix(extension, ".") {
			extension = "." + extension
		}
		name = name + extension
	}
	fileName = filepath.Join(fileName, name)
	return
}

func download(url string, fileName string) (size int64, err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := check.Get(url) // add a filter to check redirect
	if err != nil {
		return
	}

	defer resp.Body.Close()

	size, err = io.Copy(file, resp.Body)
	if err != nil {
		return
	}

	return
}

func CleanupDownloads() {
	loadConfig()
	err := filepath.Walk(settings.DownloadDir, cleanUp)
	if err != nil {
		log.Error(err)
	}
}

func cleanUp(path string, f os.FileInfo, err error) error {
	// 2 hours ago
	now, _ := strconv.Atoi(time.Now().Add(-2 * time.Hour).Format("2006010215"))

	if f.IsDir() {
		if name, err := strconv.Atoi(f.Name()); err == nil {
			if name < now {
				os.Remove(path)
			}
		}
	}
	return nil
}
