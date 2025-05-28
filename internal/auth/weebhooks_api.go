package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	if val := headers.Get("Authorization"); val != "" {
		return strings.TrimPrefix(val, "ApiKey "), nil
	}

	return "", errors.New("authorization does not exist in header")
}
