package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tomise20/go-vue/handler"
)

func main() {
	log.Println("Starting server...")

	router := gin.Default()

	handler := NewHandler(&handler.Config{
		R: router,
	})

	srv := &http.Server{
		Addr: ":8080",
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialized server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %s\n", srv.Addr)

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutfown: %v\n", err)
	}
}