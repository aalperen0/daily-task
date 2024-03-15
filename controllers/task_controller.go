package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aalperen0/dailytask/db"
	"github.com/aalperen0/dailytask/handler"
	"github.com/aalperen0/dailytask/internal/database"
	"github.com/aalperen0/dailytask/models"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func CreateTaskHandler(apiCfg *db.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			DueDate     string `json:"due_date"`
			Status      string `json:"status"`
			Email       string `json:"email"`
		}

		params := parameters{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing json", err))
			return
		}

		if params.Title == "" && params.Description == "" && params.Status == "" && params.DueDate == "" {
			handler.RespondWithError(w, http.StatusBadRequest, "All fields cannot be empty")
			return
		} else if params.Title == "" {
			handler.RespondWithError(w, http.StatusBadRequest, "Title cannot be empty")
			return
		} else if params.Description == "" {
			handler.RespondWithError(w, http.StatusBadRequest, "Description cannot be empty")
			return
		} else if params.Status == "" {
			handler.RespondWithError(w, http.StatusBadRequest, "Status cannot be empty")
			return
		} else if params.DueDate == "" {
			handler.RespondWithError(w, http.StatusBadRequest, "Due date cannot be empty")
			return
		}

		userID, err := apiCfg.DB.GetIdOfUser(r.Context(), params.Email)
		if err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Failed to get ID of user")
			return
		}

		createdAt := sql.NullTime{Time: time.Now().UTC(), Valid: true}
		updatedAt := sql.NullTime{Time: time.Now().UTC(), Valid: true}
		status := sql.NullString{String: params.Status, Valid: params.Status != ""}

		dueTime, err := time.Parse("02-01-2006", params.DueDate)
		if dueTime.String() == "" && err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Error parsing due date")
			return
		}

		dueDate := sql.NullTime{Time: dueTime, Valid: true}

		user, err := apiCfg.DB.CreateTasks(r.Context(), database.CreateTasksParams{
			ID:          uuid.New(),
			UserID:      userID,
			Title:       params.Title,
			Description: params.Description,
			DueDate:     dueDate,
			Status:      status,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})

		if err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Can not create tasks")
			return
		}

		handler.RespondWithJSON(w, http.StatusOK, models.DatabaseTaskToTask(user))
	}
}

func DeleteSpecificTaskOfUserHandler(apiCfg *db.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			TaskID uuid.UUID `json:"task_id"`
		}

		params := parameters{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse json")
			return
		}

		if params.TaskID == uuid.Nil && params.TaskID.String() == "" {
			handler.RespondWithJSON(w, http.StatusBadRequest, "TaskId doesn't exists or empty, Failed.")
			return
		}

		err := apiCfg.DB.DeleteTaskByTaskId(r.Context(), params.TaskID)
		if err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Failed to delete task")
			return
		}

		handler.RespondWithJSON(w, http.StatusOK, "Task deleted successfully")

	}
}

func DeleteAllTasksOfUserHandler(apiCfg *db.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			ID uuid.UUID `json:"id"`
		}

		params := parameters{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse json")
			return
		}

		if params.ID == uuid.Nil && params.ID.String() == "" {
			handler.RespondWithJSON(w, http.StatusBadRequest, "Id doesn't exists or empty, Failed.")
			return
		}

		err := apiCfg.DB.DeleteTaskOfUser(r.Context(), params.ID)
		if err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Failed to delete task")
			return
		}

		handler.RespondWithJSON(w, http.StatusOK, "Task deleted successfully")
	}
}

func GetTasksOfUserHandler(apiCfg *db.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			UserID uuid.UUID `json:"user_id"`
		}

		params := parameters{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse json")
			return
		}

		if params.UserID == uuid.Nil {
			handler.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		task, err := apiCfg.DB.GetTaskOfUser(r.Context(), params.UserID)
		if err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Can not get task of user")
			return
		}

		handler.RespondWithJSON(w, http.StatusOK, models.DatabaseTaskAttributes(task))
	}
}

func UpdateTaskOfUserHandler(apiCfg *db.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			TaskID      uuid.UUID `json:"task_id"`
			Description string    `json:"description"`
			Title       string    `json:"title"`
			DueDate     string    `json:"due_date"`
			Status      string    `json:"status"`
		}

		params := parameters{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse json")
			return
		}

		fmt.Printf("TaskID: %s, Description: %v, Title: %v, DueDate: %v, Status: %v\n", params.TaskID, params.Description, params.Title, params.DueDate, params.Status)

		if params.TaskID == uuid.Nil {
			handler.RespondWithError(w, http.StatusNotFound, "Task not found")
			return
		}

		existingTask, err := apiCfg.DB.GetTaskById(r.Context(), params.TaskID)
		if err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Failed to get task")
			return
		}

		var dueDate sql.NullTime
		if params.DueDate == "" {
			dueDate = existingTask.DueDate
		} else {
			dueTime, err := time.Parse("02-01-2006", params.DueDate)
			if err != nil {
				handler.RespondWithError(w, http.StatusBadRequest, "Error parsing due date")
				return
			}
			dueDate = sql.NullTime{Time: dueTime, Valid: true}
		}

		var status sql.NullString
		if params.Status == "" {
			status = existingTask.Status
		} else {
			status = sql.NullString{String: params.Status, Valid: true}
		}

		var description string
		if params.Description == "" {
			description = existingTask.Description
		} else {
			description = params.Description
		}

		var title string
		if params.Title == "" {
			title = existingTask.Title
		} else {
			title = params.Title
		}

		updatedAt := sql.NullTime{Time: time.Now().UTC(), Valid: true}

		err = apiCfg.DB.UpdateTaskOfUser(r.Context(), database.UpdateTaskOfUserParams{
			ID:          params.TaskID,
			Description: description,
			Title:       title,
			DueDate:     dueDate,
			Status:      status,
			UpdatedAt:   updatedAt,
		})
		if err != nil {
			handler.RespondWithError(w, http.StatusInternalServerError, "Failed to update tasks")
			return
		}

		handler.RespondWithJSON(w, http.StatusOK, "User's task updated successfully")
	}
}
