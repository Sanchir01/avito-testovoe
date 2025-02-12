package user

import (
	"errors"
	"github.com/Sanchir01/avito-testovoe/pkg/lib/api"
	sl "github.com/Sanchir01/avito-testovoe/pkg/lib/log"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
type AuthResponse struct {
	api.Response
	Token string `json:"token"`
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
	log.Info("sucess auth")

	render.JSON(w, r, AuthResponse{
		Response: api.OK(),
		Token:    token,
	})
}
