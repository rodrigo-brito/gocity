package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rodrigo-brito/gocity/model"

	"github.com/go-chi/chi"

	"github.com/rodrigo-brito/gocity/analyzer"
)

func main() {
	router := chi.NewRouter()
	router.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		projectName := r.URL.Query().Get("q")
		if len(projectName) == 0 {
			return
		}

		// TODO: improve folder selection, ignore GOPATH
		analyzer := analyzer.NewAnalyzer(projectName, analyzer.WithIgnoreList("/vendor/"))
		err := analyzer.FetchPackage()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Print(err)
		}

		summary, err := analyzer.Analyze()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Printf("error on analyzetion %s", err)
		}

		body, err := json.Marshal(model.New(summary, projectName))
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Print(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "frontend")
	FileServer(router, "/", http.Dir(filesDir))

	fmt.Println("Server started at http://localhost:3000")
	if err := http.ListenAndServe(":3000", router); err != nil {
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
