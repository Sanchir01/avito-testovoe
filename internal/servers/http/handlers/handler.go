package httphandlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartHTTTPHandlers(_ context.Context) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Recoverer)

	router.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware)
			r.Get("/coin", func(w http.ResponseWriter, _ *http.Request) {
				if _, err := w.Write([]byte("Hello, World!")); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

			})
		})
		r.Post("/post", func(w http.ResponseWriter, _ *http.Request) {
			if _, err := w.Write([]byte("Hello, World!")); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		})
	})

	return router

}
