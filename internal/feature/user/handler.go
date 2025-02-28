package user

import (
	"context"
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
type SendCoinsRequest struct {
	Email string `json:"toUser" validate:"required,email"`
	Coins int64  `json:"amount" validate:"required"`
}
type BuyProductResponse struct {
	api.Response
	Ok string `json:"ok"`
}
type AuthResponse struct {
	api.Response
	Token string `json:"token"`
}
type AllUserCoinsInfoResponse struct {
	api.Response
	GetAllUserCoinsInfo *GetAllUserCoinsInfo
}
type Handler struct {
	Service HandlerUser
	Log     *slog.Logger
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=HandlerUser
type HandlerUser interface {
	Auth(ctx context.Context, email, password string) (string, error)
	BuyProduct(ctx context.Context, userID, productID uuid.UUID) error
	SendCoins(ctx context.Context, userID uuid.UUID, senderEmail string, amount int64) error
	GetAllUserCoinsInfo(ctx context.Context, userID uuid.UUID) (*GetAllUserCoinsInfo, error)
}

func NewHandler(s HandlerUser, lg *slog.Logger) *Handler {
	return &Handler{
		Service: s,
		Log:     lg,
	}
}

// @Summary BuyProduct
// @Security ApiKeyAuth
// @Tags user
// @Description buy product endpoint
// @Param productid  path string true "product id"
// @Accept json
// @Produce json
// @Success 200 {string}  string "ok"
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /api/buy/{productid} [get]
func (h *Handler) BuyProductHandler(w http.ResponseWriter, req *http.Request) {
	const op = "handlers.buyProduct"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(req.Context())),
	)

	productID := chi.URLParam(req, "item")
	productuuid, err := uuid.Parse(productID)
	if err != nil {
		log.Error("failed to parse product uuid", sl.Err(err))
		render.JSON(w, req, api.Error("failed buy product"))
		return
	}

	claims, ok := req.Context().Value(contextkey.UserIDCtxKey).(*Claims)

	if !ok {
		log.Error("failed to parse product uuid")
		render.JSON(w, req, api.Error("failed buy product"))
		return
	}

	err = h.Service.BuyProduct(req.Context(), claims.ID, productuuid)
	if errors.Is(err, api.ErrInsufficientCoins) {
		log.Error("dont have coin", sl.Err(err))
		render.JSON(w, req, api.Error("недостаточно coin на балансе"))
		return
	}
	if err != nil {
		log.Error("failed to buy product", sl.Err(err))
		render.JSON(w, req, api.Error("failed, buy product"))
		return
	}

	render.JSON(w, req, BuyProductResponse{
		Response: api.OK(),
		Ok:       "success",
	})
}

// @Summary Auth
// @Tags user
// @Description buy product endpoin
// @Accept json
// @Produce json
// @Param input body AuthRequest true "auth body"
// @Success 200 {object}  AuthResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /api/auth [post]
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
	if errors.Is(err, api.ErrWrongPasswordError) {
		log.Info("password error", sl.Err(err))
		render.JSON(w, r, api.Error("Введен неправильный пароль"))
		return
	}
	if err != nil {
		log.Error("failed auth user", sl.Err(err))
		render.JSON(w, r, api.Error("failed, auth user"))
		return
	}
	log.Info("success auth")

	render.JSON(w, r, AuthResponse{
		Response: api.OK(),
		Token:    token,
	})
}

// @Summary SendUserCoins
// @Security ApiKeyAuth
// @Tags user
// @Description send coins
// @Accept json
// @Produce json
// @Param input body SendCoinsRequest true "send coins body"
// @Success 200 {string}  string "ok"
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /api/sendCoin [post]
func (h *Handler) SendUserCoinsHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.sendUserCoins"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req SendCoinsRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", sl.Err(err))
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

	err := h.Service.SendCoins(r.Context(), claims.ID, req.Email, req.Coins)
	if errors.Is(err, api.ErrTransactionMyself) {
		log.Error("send coin", slog.Any("err", err.Error()))
		render.JSON(w, r, api.Error("нельзя отправлять денеьги самому себе"))
		return
	}
	if err != nil {
		log.Error("failed send coins", sl.Err(err))
		render.JSON(w, r, api.Error("failed send coins"))
		return
	}
	log.Info("success send coins")

	render.JSON(w, r, api.OK())
}

// @Summary GetAllInfoUser
// @Security ApiKeyAuth
// @Tags user
// @Description user info coins and product buy count
// @Accept json
// @Produce json
// @Success 200 {object}  AllUserCoinsInfoResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /api/info [get]
func (h *Handler) GetInfoCoinsHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.GetUsersInfoCoins"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	claims, ok := r.Context().Value(contextkey.UserIDCtxKey).(*Claims)
	if !ok {
		log.Error("failed to parse user uuid")
		render.JSON(w, r, api.Error("failed get user coins info"))
		return
	}
	usersInfo, err := h.Service.GetAllUserCoinsInfo(r.Context(), claims.ID)
	if err != nil {
		log.Error("failed get user coins info", sl.Err(err))
		render.JSON(w, r, api.Error("failed get user coins info"))
		return
	}
	render.JSON(w, r, AllUserCoinsInfoResponse{
		Response:            api.OK(),
		GetAllUserCoinsInfo: usersInfo,
	})
}
