package handle

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rodrigo-brito/gocity/utils"

	"github.com/rodrigo-brito/gocity/analyzer"
	"github.com/rodrigo-brito/gocity/lib"
	"github.com/rodrigo-brito/gocity/model"
)

type AnalyzerHandle struct {
	Storage lib.Storage
	Cache   lib.Cache
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
