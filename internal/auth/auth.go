package auth

import (
	"net/http"
	"errors"
	"strings"
)
var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")

func GetApiKey(headers http.Header) (string, error){
        authHeader := headers.Get("Authorization")
		if authHeader == "" {
			return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	
  return splitAuth[1], nil
}