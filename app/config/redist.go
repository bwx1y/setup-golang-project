package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func ConnectionRedist() (*redis.Client, error) {
	dbNumber, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Println("Failed to connect to Redis")
		return nil, err
	}

	conn := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNumber,
	})

	_, err = conn.Ping(context.Background()).Result()
	if err != nil {
		log.Println("[ Application ] Failed to connect to Redis: %v", err)
		return nil, err
	} else {
		fmt.Println("[ Application ] Redis connected successfully")
	}

	return conn, nil
}
