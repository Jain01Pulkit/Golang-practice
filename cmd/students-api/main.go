package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jain01Pulkit/students-api/internal/config"
	"github.com/Jain01Pulkit/students-api/internal/http/handlers/student"
)


func main(){
	cfg := config.MustLoad()
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students",student.New() )

	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("Server Started",slog.String("address", cfg.Addr))

	done := make(chan os.Signal,1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func(){
		err :=server.ListenAndServe()
		if err != nil{
			log.Fatal("Failed to start server",err)
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error",err.Error()))
	}

	slog.Info("Server shutdown successfully")

}