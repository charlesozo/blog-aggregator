package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/charlesozo/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error parsing json")
		return
	}

	fmt.Println(params.Name)

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create user, %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))

}
func (cfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetPostsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get posts %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
