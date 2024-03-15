package main

import (
	"github.com/aalperen0/dailytask/controllers"
	db2 "github.com/aalperen0/dailytask/db"
	"github.com/aalperen0/dailytask/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Getting error while loading .env file ")
	}
	// GETTING PORT
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error while getting port or port is wrong")
	}

	// DATABASE CONNECTION

	apiCfg := db2.InitAPIConfig()

	// SERVER CONNECTION
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// ROUTES
	routerV1 := chi.NewRouter()
	routerV1.Get("/health", handler.HandlerReadiness)
	routerV1.Get("/error", handler.HandlerError)

	routerV1.Get("/users/{userID}", controllers.GetUserHandler(apiCfg))
	routerV1.Delete("/users/{userID}", controllers.DeleteUserHandler(apiCfg))
	routerV1.Put("/users/profile", controllers.UpdateUserHandler(apiCfg))

	routerV1.Post("/tasks", controllers.CreateTaskHandler(apiCfg))
	routerV1.Delete("/tasks/delete", controllers.DeleteAllTasksOfUserHandler(apiCfg))
	routerV1.Delete("/tasks", controllers.DeleteSpecificTaskOfUserHandler(apiCfg))
	routerV1.Get("/tasks", controllers.GetTasksOfUserHandler(apiCfg))
	routerV1.Put("/tasks/update", controllers.UpdateTaskOfUserHandler(apiCfg))

	routerV1.Post("/login", controllers.LoginHandler(apiCfg))
	routerV1.Post("/register", controllers.CreateUserHandler(apiCfg))

	router.Mount("/v1", routerV1)

	server := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server started on %v", port)
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Can not connect to server:", err)
	}
}
