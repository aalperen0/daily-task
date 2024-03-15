package controllers

import (
	"encoding/json"
	"fmt"
	db "github.com/aalperen0/dailytask/db"
	"github.com/aalperen0/dailytask/handler"
	"github.com/aalperen0/dailytask/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func generateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	fmt.Println(token)
	return token.SignedString(jwtKey)
}

func LoginHandler(apiCfg *db.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var auth Authentication

		if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse json")
		}
		if auth.Email == "" || auth.Password == "" {
			handler.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		}
		rows, err := apiCfg.DB.AuthenticateUser(r.Context(), database.AuthenticateUserParams{
			Email:    auth.Email,
			Password: auth.Password,
		})
		if err != nil {
			handler.RespondWithError(w, http.StatusInternalServerError, "Database error")
		}

		if len(rows) == 0 {
			handler.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		user := User{
			ID:        rows[0].ID,
			Email:     rows[0].Email,
			CreatedAt: rows[0].CreatedAt,
		}

		token, err := generateJWT(user.ID.String())
		if err != nil {
			handler.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		handler.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})

	}

}

var jwtKey = []byte("my_secret_key")
