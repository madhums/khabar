package log

import (
	"bytes"
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
	"strings"
)

type OurFormatter struct {
}

func (f *OurFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var keys []string = make([]string, 0, len(entry.Data))

	for k := range entry.Data {
		if k != "prefix" {
			keys = append(keys, k)
		}
	}

	b := &bytes.Buffer{}

	f.appendValue(b, strings.ToUpper(entry.Level.String()))
	f.appendValue(b, fileInfo(7))
	if entry.Message != "" {
		i := strings.LastIndex(entry.Message, "[") + 1
		f.appendValue(b, entry.Message[i:len(entry.Message)-i])
	}
	for _, key := range keys {
		if key != "file" {
			f.appendValue(b, entry.Data[key])
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return false
		}
	}
	return true
}

func (f *OurFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	fmt.Fprint(b, value)
	b.WriteByte(' ')
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	}
	return fmt.Sprintf("%s:%d:", file, line)
}
