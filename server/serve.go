package server

import (
   "time"
   "net/http"
   "github.com/gorilla/mux"
   "github.com/fuzzylemma/scowldb/pdb"
)


type Server struct {
   scowldb  *pdb.PostgresDB
   router   *mux.Router
}

func NewServer() *Server {
   server := Server{}
   server.scowldb = pdb.NewPostgresDB("")
   server.routes()
   return &server
}

func (s *Server) routes(){
   s.router = mux.NewRouter()
   s.AddRoutes()
}

func (s *Server) HttpServer() *http.Server {
   return &http.Server{
      Handler: s.router,
      Addr: "0.0.0.0:8888",
      WriteTimeout: 15 * time.Second,
      ReadTimeout: 15 * time.Second,
   }
}
