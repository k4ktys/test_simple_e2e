package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Data struct {
	m     *sync.Mutex
	Posts []string `json:"posts"`
}

func main() {
	data := Data{
		m:     &sync.Mutex{},
		Posts: []string{},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	mux.HandleFunc("/get_all", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}

		data.m.Lock()
		defer data.m.Unlock()

		resp, err := json.Marshal(&data)
		if err != nil {
			http.Error(w, "Internal unmarshal error", http.StatusInternalServerError)
		}

		w.Write(resp)
	})

	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}

		value := r.URL.Query().Get("data")

		data.m.Lock()
		defer data.m.Unlock()

		if value != "" {
			for _, post := range data.Posts {
				if value == post {
					http.Error(w, "Record is already exists", http.StatusBadRequest)
					return
				}
			}
		} else {
			http.Error(w, "Data is empty", http.StatusBadRequest)
			return
		}

		data.Posts = append(data.Posts, value)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		fmt.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("\nShutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v\n", err)
	} else {
		fmt.Println("Server stopped")
	}
}
