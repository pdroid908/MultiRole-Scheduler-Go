package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)


var (Client *goredis.Client; Ctx = context.Background())


func Connect() error{
	Client = goredis.NewClient(&goredis.Options{
		Addr: "127.0.0.1:6379",
		DB: 0,
	})

	_,err:= Client.Ping(Ctx).Result()
	return err
}


func RateLimit(key string)(bool, error){
	count,err:= Client.Incr(Ctx,key).Result()
	if err!=nil{
		return false, err
	}

	if count == 1{
		Client.Expire(Ctx,key,300*time.Second)
	}
	if count>5{
		return false,nil
	}

	return true,nil
}
