package auth

import (
	"github.com/patrickmn/go-cache"
	"github.com/jackc/pgx/v5/pgxpool"
	
)

type Data struct{
	DB *pgxpool.Pool
	Cache *cache.Cache
}