package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nv-root/task-manager/internal/config"
	"github.com/nv-root/task-manager/internal/database"
	"github.com/nv-root/task-manager/internal/handlers"
	"github.com/nv-root/task-manager/internal/repository"
	"github.com/nv-root/task-manager/internal/routes"
	"github.com/nv-root/task-manager/internal/services"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading env")
	}

	cfg := config.Config{
		MongoUri: os.Getenv("MONGO_URI"),
		Database: os.Getenv("DATABASE_NAME"),
		Port:     os.Getenv("PORT"),
	}

	client, err := database.Connect(cfg.MongoUri)
	if err != nil {
		log.Fatalln("Error Connecting to database:", err)
	}
	defer client.Disconnect(context.Background())

	taskRepo := repository.NewTaskRespository(client, cfg.Database)
	taskService := services.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	// router
	routes.TaskRouter(mux, taskHandler)

	fmt.Printf("Server running on port :%s...", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)
}
