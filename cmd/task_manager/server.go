package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		fmt.Printf("Server running on port :%s...\n", cfg.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln("Error starting server:", err)
		}
	}()

	// shutdown shit
	<-ctx.Done()
	fmt.Println("Shutting down server...")

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		log.Println("Server forced to shutdown:", err)
	}

	err = client.Disconnect(shutdownCtx)
	if err != nil {
		log.Println("Error closing database:", err)
	}

	fmt.Println("Shutdown complete.")
}
