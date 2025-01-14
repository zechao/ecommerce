package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zechao158/ecomm/config"
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
	router.HandleFunc("/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/register", h.HandleRegister).Methods("POST")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := ParseJSON(r, &payload); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		validationErr := err.(validator.ValidationErrors)
		WriteError(w, http.StatusBadRequest, validationErr)
		return
	}
	storedUser, err := h.store.GetUserByEmail(context.Background(), payload.Email)
	if err != nil && errors.Unwrap(err) == storage.ErrDuplicateKey {
		WriteError(w, http.StatusConflict, err)
		return
	}
	if !auth.ComparePassword(storedUser.Password, payload.Password) {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("wrong password"))
		return
	}

	authToken, err := auth.CreateJWT([]byte(config.ENVs.JWTSecret), storedUser.ID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{
		"token": authToken,
	})
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := ParseJSON(r, &payload); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err := h.store.GetUserByEmail(context.Background(), payload.Email)
	if err != nil && errors.Unwrap(err) == storage.ErrDuplicateKey {
		WriteError(w, http.StatusConflict, err)
		return
	}
	if err != nil && errors.Unwrap(err) != storage.ErrRecordNotFound {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		validationErr := err.(validator.ValidationErrors)
		WriteError(w, http.StatusBadRequest, validationErr)
		return
	}

	hashedPass, err := auth.HashPassword(payload.Password)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	user := &types.User{
		ID:        uuid.New(),
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPass,
	}
	err = h.store.CreateUser(context.Background(), user)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	WriteJSON(w, http.StatusCreated, &types.RegisterUserPayload{
		ID:        uuid.New(),
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
	})

}
