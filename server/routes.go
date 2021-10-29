package server

import (
   "net/http"
)

func (s *Server) AddRoutes() {
   getRequest := s.router.Methods(http.MethodGet).Subrouter()
   getRequest.HandleFunc("/random/{vrf}", GetRandomWord)
   getRequest.HandleFunc("/id/{id}", GetWordById)
   getRequest.HandleFunc("/word/{word}", GetIdByWord)
   getRequest.HandleFunc("/count", WordCount)
   getRequest.HandleFunc("/maxid", MaxId)
}
