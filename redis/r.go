package redis

import (
	"context"
	
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func Connect() error {
	// Ambil URL Redis dari environment variable Railway
	redisURL := os.Getenv("REDIS_URL") // Sesuaikan nama variabel jika di Railway berbeda (misal: REDIS_PUBLIC_URL / REDIS_PRIVATE_URL)
	
	// Jika menggunakan format URL dari Railway (contoh: redis://default:password@host:port)
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}

	Client = redis.NewClient(opt)

	_, err = Client.Ping(Ctx).Result()
	return err
}

func RateLimit(key string) (bool, error) {
	count, err := Client.Incr(Ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		Client.Expire(Ctx, key, 300*time.Second)
	}
	if count > 5 {
		return false, nil
	}

	return true, nil
}

// SetC digunakan untuk menyimpan data ke Redis Cache dengan masa kedaluwarsa (duration)
func Set(key string, value interface{}, duration time.Duration) error {
	
	_= Client.Set(Ctx,key,value,duration ).Err()
	
	return nil
}

func Get(key string)string{
	a,_:=Client.Get(Ctx,key).Result()

	return a
}

func Del(key string)error{
	_= Client.Del(Ctx,key).Err()

	return nil
}