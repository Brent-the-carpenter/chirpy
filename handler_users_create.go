package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Brent-the-carpenter/chirpy/internal/auth"
	"github.com/Brent-the-carpenter/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (state *apiConfig) handlerCreateUser(res http.ResponseWriter, req *http.Request) {
	type params struct {
		Email    string `json:"email" `
		Password string `json:"password"`
	}

	parameters := params{}

	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&parameters)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		respondWithError(res,
			http.StatusInternalServerError,
			"error unmarshilling data , error",
			err)
	}
	if parameters.Email == "" {
		respondWithError(res, http.StatusBadRequest, "missing email parameter", nil)
	}
	if parameters.Password == "" {
		respondWithError(res, http.StatusBadRequest, "Password Required", nil)
	}
	hashedPassword := auth.HashPassword(parameters.Password)

	newUser, err := state.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          parameters.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error saving user to database. error: %v", err)
		respondWithError(res, http.StatusInternalServerError, "Error saving user", err)
	}
	respondWithJSON(res, http.StatusCreated, User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	})
}
