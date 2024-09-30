package store

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SongsStore interface {
	ApplyMigrations() error
	CreateSong(group, song string) error
	GetSongs(group, song string, page, pageSize int) ([]Song, error)
	DeleteSong(songID uint64) error
	UpdateSong(songID uint64, releasedate, link, text string) error
}

type Store struct {
	db *pgxpool.Pool
}

type Song struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Link        string `json:"link"`
	Text        string `json:"text"`
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{
		db: pool,
	}
}

var (
	songTable = goqu.Dialect("postgres").From("songs").Prepared(true)
	songCols  = []any{
		"group",
		"song",
		"releasedate",
		"link",
		"text",
	}
	//go:embed 1_init.sql
	migrationFS []byte
)

var _ SongsStore = (*Store)(nil)

func (s *Store) ApplyMigrations() error {
	if _, err := s.db.Exec(context.Background(), string(migrationFS)); err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateSong(group, song string) error {
	sql, args, err := songTable.
		Insert().
		Rows(goqu.Record{
			"group": group,
			"song":  song,
		}).
		Returning(songCols...).ToSQL()
	if err != nil {
		return fmt.Errorf("create_song: %w", err)
	}

	if _, err = s.db.Exec(context.Background(), sql, args...); err != nil {
		return fmt.Errorf("exec create_song: %w", err)
	}

	return nil
}

func (s *Store) GetSongs(group, song string, page, pageSize int) ([]Song, error) {
	ds := songTable.Select(songCols...)
	if group != "" {
		ds = ds.Where(goqu.Ex{"group": group})
	}
	if song != "" {
		ds = ds.Where(goqu.Ex{"song": song})
	}
	ds = ds.Limit(uint(pageSize)).Offset(uint((page - 1) * pageSize))

	sql, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("get_songs: %w", err)
	}

	rows, err := s.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("exec get_songs: %w", err)
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song
		if err := rows.Scan(
			&song.Group,
			&song.Song,
			&song.ReleaseDate,
			&song.Link,
			&song.Text,
		); err != nil {
			return nil, fmt.Errorf("scan get_songs: %w", err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (s *Store) DeleteSong(songID uint64) error {
	sql, args, err := songTable.Where(goqu.Ex{"id": songID}).Delete().ToSQL()
	if err != nil {
		return fmt.Errorf("delete_song: %w", err)
	}

	if _, err = s.db.Exec(context.Background(), sql, args...); err != nil {
		return fmt.Errorf("exec delete_song: %w", err)
	}

	return nil
}

func (s *Store) UpdateSong(songID uint64, releasedate, link, text string) error {
	sql, args, err := goqu.Update("songs").
		Set(goqu.Record{
			"releasedate": releasedate,
			"link":        link,
			"text":        text,
		}).
		Where(goqu.Ex{"id": songID}).
		Returning(songCols...).
		ToSQL()
	if err != nil {
		return fmt.Errorf("update_song: %w", err)
	}

	if _, err = s.db.Exec(context.Background(), sql, args...); err != nil {
		return fmt.Errorf("exec update_song: %w", err)
	}

	return nil
}
