package user

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/zechao158/ecomm/service/auth"
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
	router.HandleFunc("/register", h.HandleLogin).Methods("POST")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPlayload
	if err := ParseJSON(r, payload); err != nil {
		WriteError(w, http.StatusBadRequest, err)
	}

	res, err := h.store.GetUserByEmail(context.Background(), payload.Email)
	if err != nil {
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

	err = h.store.CreateUser(context.Background(), &types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPass,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	WriteJSON(w, http.StatusCreated, res)

}
