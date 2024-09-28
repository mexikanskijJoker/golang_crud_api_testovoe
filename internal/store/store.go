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
	ApplyMigrations(ctx context.Context) error
	CreateSong(group, song string) error
	GetSongs(ctx context.Context, group, song string, page, pageSize int) ([]Song, error)
	DeleteSong(ctx context.Context, songID int) error
	UpdateSong(ctx context.Context, songID int, group, song string) error
}

type Store struct {
	db *pgxpool.Pool
}

type Song struct {
	Group string `json:"group"`
	Song  string `json:"song"`
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
	}
	//go:embed 1_init.sql
	migrationFS []byte
)

var _ SongsStore = (*Store)(nil)

func (s *Store) ApplyMigrations(ctx context.Context) error {
	if _, err := s.db.Exec(ctx, string(migrationFS)); err != nil {
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

func (s *Store) GetSongs(ctx context.Context, group, song string, page, pageSize int) ([]Song, error) {
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

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("exec get_songs: %w", err)
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.Group, &song.Song); err != nil {
			return nil, fmt.Errorf("scan get_songs: %w", err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (s *Store) DeleteSong(ctx context.Context, songID int) error {
	sql, args, err := songTable.Where(goqu.Ex{"id": songID}).Delete().ToSQL()
	if err != nil {
		return fmt.Errorf("delete_song: %w", err)
	}

	if _, err = s.db.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("exec delete_song: %w", err)
	}

	return nil
}

func (s *Store) UpdateSong(ctx context.Context, songID int, group, song string) error {
	sql, args, err := goqu.Update("songs").
		Set(goqu.Record{
			"group": group,
			"song":  song,
		}).
		Where(goqu.Ex{"id": songID}).
		Returning(songCols...).
		ToSQL()
	if err != nil {
		return fmt.Errorf("update_song: %w", err)
	}

	if _, err = s.db.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("exec update_song: %w", err)
	}

	return nil
}
