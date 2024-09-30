package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mexikanskijjoker/songs_library/internal/store"
	"github.com/sirupsen/logrus"
)

type Server struct {
	*http.Server
	store  store.SongsStore
	logger *logrus.Logger
}

func New(store store.SongsStore, logger *logrus.Logger) *Server {
	s := &Server{
		store:  store,
		logger: logger,
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

	s.logger.Infof(
		"Received GET Request, IP: %s, URL: %s, params: %v",
		r.RemoteAddr,
		r.RequestURI,
		fmt.Sprintf(
			"{ group: %s, song: %s, pageStr: %s, pageSizeStr: %s }",
			group, song, pageStr, pageSizeStr,
		),
	)

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		s.logger.Errorf("exec handleGetSongs strconv.Atoi(pageStr): %v", err)
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		s.logger.Errorf("exec handleGetSongs strconv.Atoi(pageSize): %v", err)
		pageSize = 10
	}

	songs, err := s.store.GetSongs(group, song, page, pageSize)
	if err != nil {
		s.logger.Errorf("exec handleGetSongs GetSongs(): %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		s.logger.Errorf("exec handleGetSongs json.NewEncoder.Encode(songs): %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleCreateSong(w http.ResponseWriter, r *http.Request) {
	var song store.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		s.logger.Errorf("exec handleCreateSong json.NewDecoder.Decode(&song): %v", err)
		http.Error(w, "exec handleCreateSong invalid request payload", http.StatusBadRequest)
		return
	}

	s.logger.Infof(
		"Received POST Request, IP: %s, URL: %s, params: %s",
		r.RemoteAddr,
		r.RequestURI,
		fmt.Sprintf("{ group: %s, song: %s }", song.Group, song.Song),
	)

	if err := s.store.CreateSong(song.Group, song.Song); err != nil {
		s.logger.Errorf("exec handleCreateSong CreateSong(): %v", err)
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
		s.logger.Errorf(
			"exec handleUpdateSong json.NewDecoder(r.Body).Decode(&request): %v",
			err,
		)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	s.logger.Infof(
		"Received PUT Request, IP: %s, URL: %s, params: %s",
		r.RemoteAddr,
		r.RequestURI,
		fmt.Sprintf("{ song_id: %d, releasedate: %s, link: %s, text: %s }",
			request.SongID,
			request.Song.ReleaseDate,
			request.Song.Link,
			request.Song.Text,
		),
	)

	if err := s.store.UpdateSong(
		request.SongID,
		request.Song.ReleaseDate,
		request.Song.Link,
		request.Song.Text,
	); err != nil {
		s.logger.Errorf(
			"exec handleUpdateSong UpdateSong: %v",
			err,
		)
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
		s.logger.Errorf(
			"exec handleDestroySong json.NewDecoder.Decode(&request): %v",
			err,
		)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	s.logger.Infof(
		"Received DELETE Request, IP: %s, URL: %s, params: %s",
		r.RemoteAddr,
		r.RequestURI,
		fmt.Sprintf("{ song_id: %d }",
			request.SongID,
		),
	)

	if err := s.store.DeleteSong(request.SongID); err != nil {
		s.logger.Errorf(
			"exec handleDestroySong DeleteSong: %v",
			err,
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
