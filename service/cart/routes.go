package cart

import (
	"net/http"

	"github.com/gorilla/mux"
	httputil "github.com/zechao158/ecomm/http"
	"github.com/zechao158/ecomm/types"
)

type Handler struct {
	orderStore types.OrderRepository
	orderItem  types.OrderItemRepository
}

func NewHandler(store types.CartRepository) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/carts", h.handlerlistCarts).Methods("GET")
}

func (h *Handler) handlerlistCarts(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetCarts(r.Context())
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, ps)
}
