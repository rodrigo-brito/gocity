//go:generate go run github.com/markbates/pkger/cmd/pkger -o handle

package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rodrigo-brito/gocity/analyzer"
	"github.com/rodrigo-brito/gocity/handle/middlewares"
	"github.com/rodrigo-brito/gocity/lib"
	"github.com/rodrigo-brito/gocity/model"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/markbates/pkger"
	log "github.com/sirupsen/logrus"
)

type AnalyzerHandle struct {
	Storage    lib.Storage
	Cache      lib.Cache
	projectURL *string
}

func (h *AnalyzerHandle) Handler(w http.ResponseWriter, r *http.Request) {
	projectURL, ok := lib.GetGithubBaseURL(r.URL.Query().Get("q"))
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var branch string
	if branch = r.URL.Query().Get("b"); branch == "" {
		branch = "master"
	}
	key := fmt.Sprintf("%s:%s", projectURL, branch)
	result, err := h.Cache.GetSet(key, func() ([]byte, error) {
		ok, data, err := h.Storage.Get(key)
		if err != nil {
			return nil, err
		}

		if ok && len(data) > 0 {
			return data, nil
		}

		analyzer := analyzer.NewAnalyzer(projectURL, branch, analyzer.WithIgnoreList("/vendor/"))
		err = analyzer.FetchPackage()
		if err != nil {
			return nil, err
		}

		summary, err := analyzer.Analyze()
		if err != nil {
			return nil, err
		}

		body, err := json.Marshal(model.New(summary, projectURL, branch))
		if err != nil {
			return nil, err
		}

		go func() {
			if err := h.Storage.Save(key, body); err != nil {
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
	_, _ = w.Write(result)
}

func (h *AnalyzerHandle) SetProjectURL(URL string) {
	h.projectURL = &URL
}

func (h *AnalyzerHandle) Serve(port int) error {
	router := chi.NewRouter()
	cors := middlewares.GetCors("*")
	baseURL := fmt.Sprintf("http://localhost:%d", port)

	router.Use(cors.Handler)
	router.Use(middlewares.APIHeader(fmt.Sprintf("%s/api", baseURL)))
	router.Use(middleware.DefaultLogger)

	assets, err := pkger.Open("/handle/assets")
	if err != nil {
		return err
	}

	fs := http.FileServer(assets)
	router.Handle("/*", fs)
	router.Get("/api", h.Handler)
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if h.projectURL != nil {
		log.Infof("Visualization available at: %s/#/%s", baseURL, *h.projectURL)
	} else {
		log.Infof("Server started at %s", baseURL)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
