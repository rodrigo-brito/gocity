package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/rodrigo-brito/gocity/analyzer"
	"github.com/rodrigo-brito/gocity/lib"
	"github.com/rodrigo-brito/gocity/model"
)

func main() {
	router := chi.NewRouter()
	cache := lib.NewCache()
	storage, err := lib.NewGCS(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(cors.Handler)

	router.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		projectName := r.URL.Query().Get("q")
		if len(projectName) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		result, err := cache.GetSet(projectName, func() ([]byte, error) {
			ok, data, err := storage.Get(projectName)
			if err != nil {
				return nil, err
			}

			if ok && len(data) > 0 {
				return data, nil
			}

			analyzer := analyzer.NewAnalyzer(projectName, analyzer.WithIgnoreList("/vendor/"))
			err = analyzer.FetchPackage()
			if err != nil {
				return nil, err
			}

			summary, err := analyzer.Analyze()
			if err != nil {
				return nil, err
			}

			body, err := json.Marshal(model.New(summary, projectName))
			if err != nil {
				return nil, err
			}

			// store result on Google Cloud Storage
			go func() {
				if err := storage.Save(projectName, body); err != nil {
					log.Print(err)
				}
			}()

			return body, nil
		}, time.Hour*48)

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Print(err)
			return
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	})

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "ui/build")
	FileServer(router, "/", http.Dir(filesDir))

	fmt.Println("Server started at http://localhost:4000")
	if err := http.ListenAndServe(":4000", router); err != nil {
		log.Print(err)
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
