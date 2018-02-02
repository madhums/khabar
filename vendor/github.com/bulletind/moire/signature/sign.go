package signature

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

func getValues(uri *url.URL) url.Values {
	return uri.Query()
}

func isTimestampValid(signed_on string, expiry int64) error {
	timestamp, err := time.Parse(time.RFC3339, signed_on)
	if err != nil {
		return err
	}

	current_time := time.Now()

	max_time_skew := current_time.Add(5 * time.Minute)
	if expiry == 0 {
		expiry = 60
	}
	max_time_offset := current_time.Add(time.Duration(int64(time.Minute) * -expiry))

	if timestamp.Sub(max_time_skew) > 0 {
		err := "Timestamp max skew validation error"
		log.Warn(err)
		return errors.New(err)
	}

	if timestamp.Sub(max_time_offset) < 0 {
		err := "Timestamp max offset validation error"
		log.Warn(err)
		return errors.New(err)
	}

	return nil
}

func canonicalQuery(public_key, timestamp string, expiry int64, forDelete bool) string {
	values := url.Values{"public_key": {public_key}, "timestamp": {timestamp}}
	if expiry != 0 {
		values = url.Values{
			"public_key": {public_key},
			"timestamp":  {timestamp},
			"expiry":     {strconv.FormatInt(expiry, 10)},
		}
	}
	if forDelete {
		values = url.Values{
			"public_key": {public_key},
			"timestamp":  {timestamp},
			"fordelete":  {strconv.FormatBool(true)},
		}
	}
	sorted := values.Encode()
	escaped := strings.Replace(sorted, "+", "%20", -1)
	return escaped
}

func canonicalPath(uri *url.URL) string {
	path := uri.Opaque
	if path != "" {
		path = "/" + strings.Join(strings.Split(path, "/")[3:], "/")
	} else {
		path = uri.Path
	}

	if path == "" {
		path = "/"
	}

	return path
}

func makeHmac512(message, secret string) []byte {
	key := []byte(secret)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(message))
	return h.Sum(nil)
}

func makeBase64(message []byte) string {
	encoded := base64.StdEncoding.EncodeToString(message)
	return encoded
}

func stringToSign(path, query string) string {
	val := path + "\n" + query
	return val
}

func MakeSignature(public_key, secret_key, timestamp string, expiry int64, forDelete bool, path string) string {
	//Stage1: Find public Key

	//Construct Canonical Query
	query := canonicalQuery(public_key, timestamp, expiry, forDelete)
	//Sign the strings, by joining \n
	signed_string := stringToSign(path, query)

	//Create Sha512 HMAC string
	hmac_string := makeHmac512(signed_string, secret_key)

	//Encode the resultant to base64
	base64_string := makeBase64(hmac_string)

	return base64_string
}

func IsRequestValid(
	public_key, private_key, timestamp, signature string, expiry int64, forDelete bool, path string,
) error {

	err := isTimestampValid(timestamp, expiry)
	if err != nil {
		return err
	}

	computed_signature := MakeSignature(public_key, private_key, timestamp, expiry, forDelete, path)

	if signature != computed_signature {
		return errors.New("Invalid signature")
	}

	return nil
}

func MakeUrl(public_key, secret_key string, forDelete bool, path string) string {
	subpath := path
	if strings.HasPrefix(path, "http") {
		parsed, _ := url.Parse(path)
		subpath = parsed.Path
	}

	timestamp := time.Now().Format(time.RFC3339)

	sign := MakeSignature(public_key, secret_key, timestamp, 0, forDelete, subpath)

	values := url.Values{
		"signature":  {sign},
		"timestamp":  {timestamp},
		"public_key": {public_key},
	}
	if forDelete {
		values = url.Values{
			"signature":  {sign},
			"timestamp":  {timestamp},
			"public_key": {public_key},
			"fordelete":  {strconv.FormatBool(forDelete)},
		}
	}

	sorted := values.Encode()
	escaped := strings.Replace(sorted, "+", "%20", -1)
	return path + "?" + escaped
}
