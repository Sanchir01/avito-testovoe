package user

import (
	"errors"
	"log/slog"
	"net/http"

	contextkey "github.com/Sanchir01/avito-testovoe/internal/context"

	"github.com/Sanchir01/avito-testovoe/pkg/lib/api"
	sl "github.com/Sanchir01/avito-testovoe/pkg/lib/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
type BuyProductResponse struct {
	api.Response
	Ok string `json:"ok"`
}
type AuthResponse struct {
	api.Response
	Token string `json:"token"`
}
type SendCoinsRequest struct {
	Email string `json:"toUser" validate:"required,email"`
	Coins int64  `json:"amount" validate:"required"`
}
type Handler struct {
	Service *Service
	Log     *slog.Logger
}

func NewHandler(s *Service, lg *slog.Logger) *Handler {
	return &Handler{
		Service: s,
		Log:     lg,
	}
}

func (h *Handler) BuyProductHandler(w http.ResponseWriter, req *http.Request) {
	const op = "handlers.buyProduct"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(req.Context())),
	)

	productID := chi.URLParam(req, "item")
	productuuid, err := uuid.Parse(productID)
	if err != nil {
		log.Error("failed to parse product uuid", slog.String("productID", productID))
		render.JSON(w, req, api.Error("failed buy product"))
		return
	}

	claims, ok := req.Context().Value(contextkey.UserIDCtxKey).(*Claims)

	if !ok {
		log.Error("failed to parse product uuid", slog.String("productID", productID))
		render.JSON(w, req, api.Error("failed buy product"))
		return
	}

	log.Info("attribute", slog.Any("userId", claims.ID), slog.Any("productID", productuuid))
	if err := h.Service.BuyProduct(req.Context(), claims.ID, productuuid); err != nil {
		log.Error("failed to buy product", sl.Err(err))
		render.JSON(w, req, api.Error("failed, buy product"))
		return
	}
	render.JSON(w, req, BuyProductResponse{
		Response: api.OK(),
		Ok:       "success",
	})
}

func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req AuthRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.JSON(w, r, api.Error("Ошибка при валидации данных"))
		return
	}
	log.Info("request body decoded", slog.Any("request", req))
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	token, err := h.Service.Auth(r.Context(), req.Email, req.Password)
	if errors.Is(err, errors.New("Неправльный пароль")) {
		log.Info("password error", slog.String("password", req.Password))
		render.JSON(w, r, api.Error("Введен неправильный пароль"))
		return
	}
	if err != nil {
		log.Error("failder auth user", sl.Err(err))
		render.JSON(w, r, api.Error("failed, auth user"))
		return
	}
	log.Info("success auth")

	render.JSON(w, r, AuthResponse{
		Response: api.OK(),
		Token:    token,
	})
}

func (h *Handler) SendUserCoinsHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.sendUserCoins"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req SendCoinsRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.JSON(w, r, api.Error("Ошибка при валидации данных"))
		return
	}
	log.Info("request body decoded", slog.Any("request", req))
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	claims, ok := r.Context().Value(contextkey.UserIDCtxKey).(*Claims)
	if !ok {
		log.Error("failed to parse product uuid")
		render.JSON(w, r, api.Error("failed send coins"))
		return
	}

	if err := h.Service.SendCoins(r.Context(), claims.ID, req.Email, req.Coins); err != nil {
		log.Error("failed send coins", sl.Err(err))
		render.JSON(w, r, api.Error("failed send coins"))
		return
	}
	log.Info("success send coins")

	render.JSON(w, r, api.OK())
}

func (h *Handler) GetInfoCoinsHandler(w http.ResponseWriter, r *http.Request) {}
