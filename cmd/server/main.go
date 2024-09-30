package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mexikanskijjoker/songs_library/internal/logger"
	"github.com/mexikanskijjoker/songs_library/internal/server"
	"github.com/mexikanskijjoker/songs_library/internal/store"
)

const (
	databaseURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx); err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fmt.Println("gracefully stopped")
}

func run(ctx context.Context) error {
	pgxCfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return err
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	storage := store.New(pool)
	if err := storage.ApplyMigrations(); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	logger := logger.New()
	s := server.New(storage, logger)

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.Shutdown(shutdownCtx); err != nil {
			log.Printf("error shutting down http server: %v", err)
		}
	}()

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
