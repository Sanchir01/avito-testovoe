package product

import (
	"context"
	"github.com/Sanchir01/avito-testovoe/pkg/lib/api"
	sl "github.com/Sanchir01/avito-testovoe/pkg/lib/log"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type Handler struct {
	Service HandlerProducts
	Log     *slog.Logger
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=HandlerProducts
type HandlerProducts interface {
	GetAllProducts(ctx context.Context) ([]*DataBaseProduct, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*DataBaseProduct, error)
}
type GetAllProductsResponse struct {
	api.Response
	Products []*DataBaseProduct `json:"products"`
}

func NewHandler(service HandlerProducts, log *slog.Logger) *Handler {
	return &Handler{
		Service: service,
		Log:     log,
	}
}

func (h *Handler) GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.Service.GetAllProducts(r.Context())
	if err != nil {
		h.Log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
	}
	render.JSON(w, r, GetAllProductsResponse{
		Response: api.OK(),
		Products: products,
	})
}
