package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	// import /setup/setup.go
	setup "veloce/setup"
)

func main() {
	config := setup.SetupAppDir()
	appdir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	appdir = filepath.Join(appdir, ".veloce")

	handler := func(w http.ResponseWriter, r *http.Request) {
		publicDir := config["public"].(string)

		for _, ignore := range config["ignore"].([]interface{}) {
			if r.URL.Path == "/"+ignore.(string) {
				http.NotFound(w, r)
				return
			}
		}

		// serve static files
		if r.URL.Path != "/" {
			http.ServeFile(w, r, filepath.Join(appdir, publicDir, r.URL.Path[1:]))
			return
		}
	}

	s := &http.Server{
		Addr:           ":" + config["port"].(string),
		Handler:        http.HandlerFunc(handler),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	fmt.Println("Server listening on port " + s.Addr)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server gracefully stopped")
}
