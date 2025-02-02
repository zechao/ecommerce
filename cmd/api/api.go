package api

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zechao158/ecomm/service/product"
	"github.com/zechao158/ecomm/service/user"
	"gorm.io/gorm"
)

type APIServer struct {
	addr string
	db   *gorm.DB
}

func NewAPIServer(addr string, db *gorm.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	router.Use(PanicRecoveryMiddleware)
	router.Path("/health").Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "ok")
		})

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewRepository(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewRepository(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	log.Println("Http servevr listening on:", s.addr)
	server := &http.Server{
		Addr:    s.addr,
		Handler: router,
	}
	return server.ListenAndServe()

}
