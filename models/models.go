package models

import (
	"fmt"
	"github.com/aalperen0/dailytask/internal/database"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func DatabaseUserToUser(dbUser database.User) User {
	var createdAt time.Time
	if dbUser.CreatedAt.Valid {
		createdAt = dbUser.CreatedAt.Time
	}
	return User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		CreatedAt: createdAt,
	}
}

type Task struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func DatabaseTaskToTask(dbTask database.Task) Task {
	var createdAt, updatedAt, dueDate time.Time
	var status string
	if dbTask.CreatedAt.Valid && dbTask.UpdatedAt.Valid && dbTask.DueDate.Valid && dbTask.Status.Valid {
		createdAt = dbTask.CreatedAt.Time
		updatedAt = dbTask.UpdatedAt.Time
		dueDate = dbTask.DueDate.Time
		status = dbTask.Status.String
	}

	return Task{
		ID:          dbTask.ID,
		UserID:      dbTask.UserID,
		Title:       dbTask.Title,
		Description: dbTask.Description,
		DueDate:     dueDate,
		Status:      status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

type TaskAttributes struct {
	Title       string
	Description string
	Status      string
	DueDate     string
}

func DatabaseTaskAttributes(tasks []database.GetTaskOfUserRow) []TaskAttributes {

	var task []TaskAttributes
	for _, t := range tasks {
		status := ""
		if t.Status.Valid {
			status = t.Status.String
		}
		dueDate := ""
		if t.DueDate.Valid {
			dueDate = t.DueDate.Time.Format("02-01-2006")
		}

		taskAttribute := TaskAttributes{
			Title:       t.Title,
			Description: t.Description,
			Status:      status,
			DueDate:     dueDate,
		}
		task = append(task, taskAttribute)
	}
	fmt.Println(task)
	return task
}
