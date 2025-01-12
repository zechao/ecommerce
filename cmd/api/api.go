package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewRepository(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)
	
	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
