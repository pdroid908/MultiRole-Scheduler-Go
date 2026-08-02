package modules

import (
	
	"github.com/jackc/pgx/v5/pgxpool"
	
)

type Data struct{
	DB *pgxpool.Pool
}