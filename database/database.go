package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("env db url kosong")
	}

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("gagal buat config %v", err)
	}
	

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("gagal buat pool %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()                                     // Perbaikan: Amankan/bersihkan pool jika ping gagal
		return nil, fmt.Errorf("gagal ping DB: %w", err) // Perbaikan: Sertakan error asli
	}
	return pool, nil
}
