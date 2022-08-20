package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/rodrigo-brito/gocity/pkg/analyzer"
	"github.com/rodrigo-brito/gocity/pkg/lib"
	"github.com/rodrigo-brito/gocity/pkg/model"
	"github.com/rodrigo-brito/gocity/pkg/server/middlewares"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

//go:embed assets
var assets embed.FS

//go:embed assets/index.html
var indexPage []byte

var ErrInvalidPath = fmt.Errorf("invalid path")

type AnalyzerHandle struct {
	Cache       lib.Cache
	CacheTTL    time.Duration
	TmpFolder   string
	Port        int
	ProjectPath *string
	Branch      *string
	Local       bool
}

func (h *AnalyzerHandle) Handler(w http.ResponseWriter, r *http.Request) {
	var (
		ok             bool
		projectAddress string
	)

	if h.ProjectPath != nil {
		projectAddress = *h.ProjectPath
	}

	if q := r.URL.Query().Get("q"); q != "local" {
		projectAddress, ok = lib.GetGithubBaseURL(r.URL.Query().Get("q"))
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	branch := "master"
	if b := r.URL.Query().Get("b"); b != "" {
		branch = b
	} else if h.Branch != nil {
		branch = *h.Branch
	}

	key := fmt.Sprintf("%s:%s", projectAddress, branch)
	result, err := h.Cache.GetSet(key, func() ([]byte, error) {
		codeAnalyzer := analyzer.NewAnalyzer(projectAddress, branch, h.TmpFolder, analyzer.WithIgnoreList("/vendor/"))

		path := projectAddress
		stat, err := os.Stat(projectAddress)
		if os.IsNotExist(err) {
			path, err = codeAnalyzer.FetchPackage()
			if err != nil {
				return nil, err
			}
		} else if !stat.IsDir() {
			return nil, ErrInvalidPath
		}

		summary, err := codeAnalyzer.Analyze(path)
		if err != nil {
			return nil, err
		}

		body, err := json.Marshal(model.New(summary, projectAddress, branch))
		if err != nil {
			return nil, err
		}

		return body, nil
	}, h.CacheTTL)

	if err == ErrInvalidPath {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
	} else if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Print(err)
		return
	}

	if len(result) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(result)
	if err != nil {
		log.Error(err)
	}
}

func (h *AnalyzerHandle) Serve() error {
	router := chi.NewRouter()
	cors := middlewares.GetCors("*")
	baseURL := fmt.Sprintf("http://localhost:%d", h.Port)

	router.Use(cors.Handler)
	router.Use(middlewares.APIHeader(fmt.Sprintf("%s/api", baseURL)))
	router.Use(middleware.DefaultLogger)

	router.Get("/api", h.Handler)
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write(indexPage)
		if err != nil {
			log.Error(err)
		}
	})

	var staticFS = fs.FS(assets)
	content, err := fs.Sub(staticFS, "assets")
	if err != nil {
		log.Fatal(err)
	}
	router.Handle("/*", http.FileServer(http.FS(content)))

	if h.Local {
		log.Infof("Visualization available at: %s/#/local", baseURL)
	} else if h.ProjectPath != nil {
		log.Infof("Visualization available at: %s/#/%s", baseURL, *h.ProjectPath)
	} else {
		log.Infof("Server started at %s", baseURL)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", h.Port), router)
}
