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
	r.HandleFunc("/api/v1/info", s.handleGetSongs).Methods("GET")
	r.HandleFunc("/api/v1/create", s.handleCreateSong).Methods("POST")
	r.HandleFunc("/api/v1/update", s.handleUpdateSong).Methods("PUT")
	r.HandleFunc("/api/v1/delete", s.handleDestroySong).Methods("DELETE")

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

	if err := s.store.CreateSong(song.Group, song.Song); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUpdateSong(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SongID uint64     `json:"song_id"`
		Song   store.Song `json:"song_detail"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := s.store.UpdateSong(
		request.SongID,
		request.Song.ReleaseDate,
		request.Song.Link,
		request.Song.Text,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleDestroySong(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SongID uint64 `json:"song_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := s.store.DeleteSong(request.SongID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
