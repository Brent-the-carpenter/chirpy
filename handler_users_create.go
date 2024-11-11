package main

import (
	"encoding/json"
	"github.com/Brent-the-carpenter/chirpy/internal/auth"
	"github.com/Brent-the-carpenter/chirpy/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
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
		return
	}
	if parameters.Email == "" || parameters.Password == "" {
		respondWithError(res, http.StatusBadRequest, "missing email or password parameter", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(parameters.Password)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "could not hash password", err)
		return
	}

	newUser, err := state.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          parameters.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error saving user to database. error: %v", err)
		respondWithError(res, http.StatusInternalServerError, "Error saving user", err)
		return
	}
	respondWithJSON(res, http.StatusCreated, User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	})
}
