package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abiiranathan/acada/database"
	"github.com/abiiranathan/acada/handlers"
	"github.com/abiiranathan/acada/services"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	flag.Parse()

	db := database.MustConnect(os.Getenv("DSN"))
	database.MigrateAll(db)
	srv := services.New(db)
	r := handlers.SetupHandlers(srv)

	IfCreateAdmin(db)
	startServer(r)
}

func startServer(handler http.Handler) {
	contextTimeout := 10 * time.Second

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   60 * time.Second,
	}

	done := make(chan error, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		ctx := context.Background()
		var cancel context.CancelFunc
		if contextTimeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, contextTimeout)
			defer cancel()
		}
		done <- server.Shutdown(ctx)
	}()

	log.Printf("http server started on http:0.0.0.0:%s\n", os.Getenv("PORT"))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped with error: %s\n", err)
	}
	log.Println("Server stopped gracefully")
}
