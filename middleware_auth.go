package main

import (
	"fmt"
	"github.com/charlesozo/blog-aggregator/internal/auth"
	"github.com/charlesozo/blog-aggregator/internal/database"
	"net/http"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) MiddleWareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid API key- %s", err))
			return
		}
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error fetching user - %s", err))
			return
		}
		handler(w, r, user)

	}
}
