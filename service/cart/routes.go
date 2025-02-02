package cart

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	httputil "github.com/zechao158/ecomm/http"
	"github.com/zechao158/ecomm/service/auth"
	"github.com/zechao158/ecomm/types"
)

type Handler struct {
	uowStore UnitOfWork
}

func NewHandler(uow UnitOfWork) *Handler {
	return &Handler{
		uowStore: uow,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/checkout", h.handlerCheckout).Methods("POST")
}

func (h *Handler) handlerCheckout(w http.ResponseWriter, r *http.Request) {
	var cart types.CartCheckoutPayload
	if err := httputil.ParseJSON(r, &cart); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := httputil.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		httputil.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	for _, eachItem := range cart.Items {
		if eachItem.Quantity <= 0 {
			httputil.WriteError(w, http.StatusBadRequest, fmt.Errorf("quantity must be greater than 0"))
			return
		}
	}

	// create the order
	// create order items

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
	}
	h.uowStore.Do(func(store OrderUOWStore) error {
		ps, err := store.productRepository.GetProductsByIDs(r.Context(), getItemsIds(cart))
		if err != nil {
			httputil.WriteError(w, http.StatusInternalServerError, err)
			return err
		}
		productMap := make(map[uuid.UUID]types.Product)
		for _, product := range ps {
			productMap[product.ID] = product
		}
		// check if all products are actually in stock
		if err := checkIfCartIsInStock(cart.Items, productMap); err != nil {
			httputil.WriteError(w, http.StatusBadRequest, err)
			return err
		}
		// calculate total price

		totalPrice := calculateTotalPrice(cart.Items, productMap)

		// reduce quantity of each product
		for _, item := range cart.Items {
			product := productMap[item.ProductID]
			product.Quantity -= item.Quantity
			if err := store.productRepository.Update(r.Context(), &product); err != nil {
				httputil.WriteError(w, http.StatusInternalServerError, err)
				return err
			}
		}

		order := types.Order{
			ID:      uuid.New(),
			UserID:  user.ID,
			Total:   totalPrice,
			Status:  "pending",
			Address: "some address",
		}
		err = store.orderRepository.Create(r.Context(), &order)
		if err != nil {
			httputil.WriteError(w, http.StatusInternalServerError, err)
			return err
		}

		httputil.WriteJSON(w, http.StatusOK, map[string]any{
			"order_id": order.ID,
			"total":    totalPrice,
		})
		return nil
	})
}

func calculateTotalPrice(cartItem []types.CartItem, productMap map[uuid.UUID]types.Product) float64 {
	var total float64
	for _, item := range cartItem {
		product := productMap[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}
	return total
}

func getItemsIds(cart types.CartCheckoutPayload) []uuid.UUID {
	ids := make([]uuid.UUID, len(cart.Items))
	for i := range cart.Items {
		ids[i] = cart.Items[i].ProductID
	}
	return ids
}

func checkIfCartIsInStock(cartItems []types.CartItem, productMap map[uuid.UUID]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}
	for _, item := range cartItems {
		product, ok := productMap[item.ProductID]
		if !ok {
			return fmt.Errorf("product not found")
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is out of stock", product.Name)
		}
	}
	return nil
}
