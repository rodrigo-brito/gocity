package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rodrigo-brito/gocity/model"

	"github.com/go-chi/chi"

	"github.com/rodrigo-brito/gocity/analyzer"
)

func main() {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		projectName := r.URL.Query().Get("q")
		w.Write([]byte(projectName))

		// TODO: validate project name

		analyzer := analyzer.NewAnalyzer(projectName)
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

		w.Write(body)
	})

	fmt.Println("Server started at http://localhost:3000")
	http.ListenAndServe(":3000", router)
}
