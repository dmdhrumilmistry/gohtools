package http

import (
	"net/url"
	"strings"
)

func IsValidURL(link string) bool {
	u, err := url.Parse(link)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Only predicts type from link, it doesn't make call to the server
func GetDocTypeFromLink(url string) string {
	tokens := strings.Split(url, ".")
	if len(tokens) < 2 {
		return "unknown"
	}
	return tokens[len(tokens)-1]
}
