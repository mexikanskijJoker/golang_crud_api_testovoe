package server

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	r.HandleFunc("/api/v1/songs/index", s.handleGetSongs).Methods("GET")
	r.HandleFunc("/api/v1/songs/create", s.handleCreateSong).Methods("POST")

	s.Server = &http.Server{
		Handler: r,
		Addr:    ":8080",
	}

	return s
}

func (s *Server) handleGetSongs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	songs, err := s.store.GetSongs(r.Context(), group, song, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleCreateSong(w http.ResponseWriter, r *http.Request) {
	var song store.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "exec handleCreateSong invalid request payload", http.StatusBadRequest)
		return
	}

	if err := store.SongsStore.CreateSong(&store.Store{}, song.Group, song.Song); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
