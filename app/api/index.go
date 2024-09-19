package api

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/deniskrumko/visual-coordinates-tool/pkg/env"
	"github.com/go-chi/chi/v5"
)

const (
	IndexTemplate      = "templates/index.html"
	DefaultDisplayName = "Visual Coordinates Tool"
)

// Context to render index page
type IndexPageContext struct {
	DisplayName   string
	ServiceGroups map[env.ServiceGroup][]env.Service
	Samples       []string
}

// Get index page router with handlers
func getIndexRouter(config *env.Config) (chi.Router, error) {
	r := chi.NewRouter()

	t, err := template.ParseFiles(IndexTemplate)
	if err != nil {
		return nil, fmt.Errorf("can't parse index template: %w", err)
	}

	tContext := IndexPageContext{
		DisplayName: DefaultDisplayName,
	}

	// If config is presented – use params from it
	if config != nil {
		if configName := config.Common.Name; configName != "" {
			tContext.DisplayName = configName
		}

		tContext.ServiceGroups = config.GetServiceGroups()
		tContext.Samples = config.GetSamples()
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, tContext); err != nil {
			errorResponse(w, fmt.Errorf("can't execute index template: %w", err))
		}
	})

	return r, nil
}
