package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mexikanskijjoker/songs_library/internal/store"
)

type Server struct {
	*http.Server
	store store.SongsStore
}

func New(store store.SongsStore) *Server {
	s := &Server{
		store: store,
	}

	r := mux.NewRouter()

	s.Server = &http.Server{
		Handler: r,
		Addr:    ":8080",
	}

	return s
}
