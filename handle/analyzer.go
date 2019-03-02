package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rodrigo-brito/gocity/handle/middlewares"

	log "github.com/sirupsen/logrus"

	"github.com/rodrigo-brito/gocity/utils"

	"github.com/rodrigo-brito/gocity/analyzer"
	"github.com/rodrigo-brito/gocity/lib"
	"github.com/rodrigo-brito/gocity/model"
)

type AnalyzerHandle struct {
	Storage    lib.Storage
	Cache      lib.Cache
	projectURL *string
}

func (h *AnalyzerHandle) Handler(w http.ResponseWriter, r *http.Request) {
	projectURL, ok := utils.GetGithubBaseURL(r.URL.Query().Get("q"))
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.Cache.GetSet(projectURL, func() ([]byte, error) {
		ok, data, err := h.Storage.Get(projectURL)
		if err != nil {
			return nil, err
		}

		if ok && len(data) > 0 {
			return data, nil
		}

		analyzer := analyzer.NewAnalyzer(projectURL, analyzer.WithIgnoreList("/vendor/"))
		err = analyzer.FetchPackage()
		if err != nil {
			return nil, err
		}

		summary, err := analyzer.Analyze()
		if err != nil {
			return nil, err
		}

		body, err := json.Marshal(model.New(summary, projectURL))
		if err != nil {
			return nil, err
		}

		// store result on Google Cloud Storage
		go func() {
			if err := h.Storage.Save(projectURL, body); err != nil {
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

	router.Get("/*", FrontEndProxy)
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
