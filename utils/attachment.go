package utils

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"path/filepath"

	"github.com/bulletind/moire/signature"
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

func DownloadFile(url string, fileName string, isPrivate bool) (filePath string, size int64, err error) {
	loadConfig()

	filePath, err = getFilePath(url, fileName)
	if err != nil {
		return
	}

	// simply check if file exists
	var fi os.FileInfo
	if fi, err = os.Stat(filePath); err == nil {
		size = fi.Size()
		return
	}

	// so we need to download
	url = MakeUrl(url, isPrivate, true, 0)
	size, err = download(url, filePath)
	return
}

func MakeUrl(url string, isPrivate bool, skip_check bool, expiry int64) string {
	if !isPrivate {
		return url
	}
	if settings.PublicKey == "" {
		log.Println("When using private urls, you need to provide keys for the mediaserver")
		return url
	}

	// add host info when needed
	if !strings.HasPrefix(url, "http") {
		if !strings.HasSuffix(settings.Host, "/") && !strings.HasPrefix(url, "/") {
			url = "/" + url
		}
		url = settings.Host + url
	}
	url = signUrl(strings.Split(url, "?")[0], expiry)
	if skip_check {
		url = url + "&skip_ready_check=true"
	}
	return url
}

func signUrl(path string, expiry int64) string {
	subpath := path
	if strings.HasPrefix(path, "http") {
		parsed, _ := url.Parse(path)
		subpath = parsed.Path
	}

	timestamp := time.Now().Format(time.RFC3339)

	sign := signature.MakeSignature(settings.PublicKey, settings.SecretKey, timestamp, expiry, false, subpath)

	values := url.Values{
		"signature":  {sign},
		"timestamp":  {timestamp},
		"public_key": {settings.PublicKey},
		"expiry":     {strconv.FormatInt(expiry, 10)},
	}

	sorted := values.Encode()
	escaped := strings.Replace(sorted, "+", "%20", -1)
	return path + "?" + escaped
}

func getFilePath(rawURL string, fileName string) (filePath string, err error) {
	fileURL, err := url.Parse(rawURL)

	if err != nil {
		return
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	// add date path so we can re-use downloaded files
	filePath = filepath.Join(settings.DownloadDir, time.Now().Format("20060102"))
	filePath = filepath.Join(filePath, fileURL.Host)
	for i := 1; i < len(segments)-1; i++ {
		filePath = filepath.Join(filePath, segments[i])
	}

	// some urls may contain filename
	// if so, we replcae that by our provided name, otherwise we add the filename
	lastSegment := segments[len(segments)-1]
	if !strings.Contains(lastSegment, ".") {
		filePath = filepath.Join(filePath, lastSegment)
		err = os.MkdirAll(filePath, 0777)
		if err != nil {
			return
		}
		reg, err := regexp.Compile("[^A-Za-z0-9.]+")
		if err != nil {
			log.Println(err)
		}
		filePath = filepath.Join(filePath, reg.ReplaceAllString(fileName, "_"))
	} else if len(fileName) == 0 {
		err = os.MkdirAll(filePath, 0777)
		if err != nil {
			return
		}
		// we have no filename, so let's use the one that's provided
		filePath = filepath.Join(filePath, lastSegment)
	}
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
		log.Println(err)
	}
}

func cleanUp(path string, f os.FileInfo, err error) error {
	// 2 days ago
	now, _ := strconv.Atoi(time.Now().Add(-48 * time.Hour).Format("20060102"))

	if f.IsDir() {
		if name, err := strconv.Atoi(f.Name()); err == nil {
			if name < now {
				log.Println("remove:", path)
				os.RemoveAll(path)
			}
		}
	}
	return nil
}
