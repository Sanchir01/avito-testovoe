package httphandlers

import (
	"net/http"

	"github.com/Sanchir01/avito-testovoe/internal/servers/http/custommiddleware"

	"github.com/Sanchir01/avito-testovoe/internal/app"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartHTTTPHandlers(handlers *app.Handlers) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Recoverer)

	router.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.AuthMiddleware)
			r.Get("/info", func(w http.ResponseWriter, _ *http.Request) {
				if _, err := w.Write([]byte("Hello, World!")); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			})
			r.Get("/buy/{item}", handlers.UserHandler.BuyProductHandler)
			r.Get("/info", handlers.UserHandler.GetInfoCoinsHandler)
			r.Post("/sendCoin", handlers.UserHandler.SendUserCoinsHandler)

		})
		r.Post("/auth", handlers.UserHandler.AuthHandler)
	})

	return router
}
