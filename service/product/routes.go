package product

import (
	"net/http"

	"github.com/gorilla/mux"
	httputil "github.com/zechao158/ecomm/http"
	"github.com/zechao158/ecomm/types"
)

type Handler struct {
	store types.ProductRepository
}

func NewHandler(store types.ProductRepository) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("", h.handlerlistProducts).Methods("GET")
}

func (h *Handler) handlerlistProducts(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetAll(r.Context(), nil)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, ps)
}
