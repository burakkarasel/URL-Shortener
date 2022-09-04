package helpers

import (
	"os"
	"strings"
)

// RemoveDomainError checks if request comes for the same URL
func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}
	newUrl := strings.Replace(url, "http://", "", 1)
	newUrl = strings.Replace(newUrl, "https://", "", 1)
	newUrl = strings.Replace(newUrl, "www.", "", 1)
	newUrl = strings.Split(newUrl, "/")[0]

	return !(newUrl == os.Getenv("DOMAIN"))
}

// EnforceHTTP we add http if it's not specified
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}
