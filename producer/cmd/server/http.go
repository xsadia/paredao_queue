package server

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	Port   string
	Router *mux.Router
}

func (s *Server) Run() {
	ma := handlers.MaxAge(60)
	ao := handlers.AllowedOrigins([]string{"*"})
	am := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"})
	ah := handlers.AllowedHeaders([]string{"Content-Type", "Origin", "Authorization"})

	log.Fatal(http.ListenAndServe(s.Port, handlers.CORS(ma, ao, am, ah)(s.Router)))
}

func NewServer(port string) Server {
	return Server{
		port,
		mux.NewRouter(),
	}
}
