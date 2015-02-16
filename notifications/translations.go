package notifications

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
)

var loadedFileCache = map[string]bool{}

func LoadTranslationFile(filename string) error {
	if loadedFileCache[filename] {
		log.Println("Already cached.")
		return nil
	}

	err := i18n.LoadTranslationFile(filename)
	if err == nil {
		loadedFileCache[filename] = true
	}
	return err
}
