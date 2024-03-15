package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aalperen0/dailytask/db"
	"github.com/aalperen0/dailytask/handler"
	"github.com/aalperen0/dailytask/internal/database"
	"github.com/aalperen0/dailytask/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func CreateUserHandler(apiCfg *db.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type parameters struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		params := parameters{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing json:", err))
			return
		}

		if params.Password == "" || params.Email == "" || params.Name == "" {
			handler.RespondWithError(w, http.StatusMethodNotAllowed, fmt.Sprintf("these fields can not be empty"))
			return
		}

		userCreatedAt := sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		}

		user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			Name:      params.Name,
			Email:     params.Email,
			Password:  params.Password,
			CreatedAt: userCreatedAt,
		})

		if err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create user: ", err))
		}

		handler.RespondWithJSON(w, http.StatusCreated, models.DatabaseUserToUser(user))
	}
}

func GetUserHandler(apiCfg *db.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDstr := chi.URLParam(r, "userID")
		log.Printf(userIDstr)
		userID, err := uuid.Parse(userIDstr)
		if err != nil {
			handler.RespondWithJSON(w, http.StatusBadRequest, "Invalid User ID")
			return
		}

		user, err := apiCfg.DB.GetUserById(r.Context(), userID)
		if err != nil {
			handler.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch user")
			return
		}

		handler.RespondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
	}
}

func DeleteUserHandler(apiCfg *db.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "userID")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse user with id:%v", err))
			return
		}
		err = apiCfg.DB.DeleteUserById(r.Context(), userID)
		if err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Failed to delete user!")
			return
		}

		handler.RespondWithJSON(w, http.StatusOK, "User deleted successfully")
	}
}

func UpdateUserHandler(apiCfg *db.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			UserID   string `json:"userID"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse JSON")
			return
		}

		userID, err := uuid.Parse(params.UserID)
		if err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		user, err := apiCfg.DB.GetUserById(r.Context(), userID)
		if err != nil {
			handler.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch user")
			return
		}

		if user.ID == uuid.Nil {
			handler.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		err = apiCfg.DB.UpdateUserById(r.Context(), database.UpdateUserByIdParams{
			Name:     params.Name,
			Email:    params.Email,
			Password: params.Password,
		})
		if err != nil {
			handler.RespondWithError(w, http.StatusInternalServerError, "Failed to update user")
			return
		}

		handler.RespondWithJSON(w, http.StatusOK, "User updated successfully")
	}
}
