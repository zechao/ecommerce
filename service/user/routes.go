package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zechao158/ecomm/config"
	httputil "github.com/zechao158/ecomm/http"
	"github.com/zechao158/ecomm/service/auth"
	"github.com/zechao158/ecomm/storage"
	"github.com/zechao158/ecomm/types"
)

type Handler struct {
	store types.UserRepository
}

func NewHandler(store types.UserRepository) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := httputil.ParseJSON(r, &payload); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := httputil.Validate.Struct(payload); err != nil {
		validationErr := err.(validator.ValidationErrors)
		httputil.WriteError(w, http.StatusBadRequest, validationErr)
		return
	}
	storedUser, err := h.store.GetUserByEmail(r.Context(), payload.Email)
	if err != nil && errors.Is(err, storage.ErrRecordNotFound) {
		httputil.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}
	if !auth.ComparePassword(storedUser.Password, payload.Password) {
		httputil.WriteError(w, http.StatusBadRequest, fmt.Errorf("wrong password"))
		return
	}

	authToken, err := auth.CreateJWT([]byte(config.ENVs.JWTSecret), storedUser.ID)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]string{
		"token": authToken,
	})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := httputil.ParseJSON(r, &payload); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := httputil.Validate.Struct(payload); err != nil {
		validationErr := err.(validator.ValidationErrors)
		httputil.WriteError(w, http.StatusBadRequest, validationErr)
		return
	}

	hashedPass, err := auth.HashPassword(payload.Password)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	user := types.User{
		ID:        uuid.New(),
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPass,
	}
	err = h.store.Create(r.Context(), &user)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateKey) {
			httputil.WriteError(w, http.StatusConflict, err)
			return
		}
		httputil.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, &types.RegisterUserPayload{
		ID:        uuid.New(),
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
	})

}
