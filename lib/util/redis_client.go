///
/// *** This file may be removed due to future architecture changes. ***
///
package util

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

// redis.Options default pool size is 10 times the number of CPUs.

const (
	Empty = ""
	FAIL  = -1
)

type RedisClient struct {
	Host string
	db   int
	pool *redis.Client
}

// NewRedisClient returns a new RedisClient.
func NewRedisClient(host string, db int) *RedisClient {
	redis_client := &RedisClient{Host: host, db: db}
	// Connection pool.
	redis_client.pool = redis.NewClient(&redis.Options{
		Addr:            host,
		DB:              0,
		MaxRetries:      3,
		MinRetryBackoff: 100 * time.Millisecond, // Minimum backoff between each retry.
		DialTimeout:     5 * time.Second,
		//WriteTimeOut:    1 * time.Second,
		PoolSize: 120,
		//MaxConnLifetime: 10 * time.Second, // Maximum age of unused connections.
		IdleTimeout: 10 * time.Second, // Maximum amount of time a connection may be idle before the connection is closed.
	})
	// Test connection.
	_, err := redis_client.pool.Ping().Result()
	if err != nil {
		log.Fatal("Redis connection error: ", err)
	}
	log.Println("Redis connection successful.")
	ping := redis_client.Ping()
	if ping == "PONG" {
		log.Println("Redis ping successful.")
	}
	return redis_client
}

func (redis_client *RedisClient) Ping() string {
	ping := redis_client.pool.Ping()
	result, err := ping.Result()
	if err != nil {
		log.Println("Error when pinging Redis: ", err)
	}
	return result
}

// Keys returns all keys with the given pattern
// in the database, in the form of a string slice.
func (redis_client *RedisClient) Keys(s string) []string {
	result, err := redis_client.pool.Keys(s).Result()
	if err != nil {
		log.Printf("Error when getting %s keys: %v", s, err)
		return nil
	}
	return result
}

// Set sets the value of a key.
func (redis_client *RedisClient) Set(key string, value string, expiration time.Duration) string {
	result, err := redis_client.pool.Set(key, value, expiration).Result()
	if err != nil {
		log.Println("Error when setting key: ", err)
		return ""
	}
	return result
}

// Get returns the value of a key.
func (redis_client *RedisClient) Get(key string) string {
	result, err := redis_client.pool.Get(key).Result()
	if err != nil {
		log.Printf("Error when getting key %s: %v", key, err)
		return Empty
	}
	return result
}

// GetSet returns the value of a key and sets it to a new value.
func (redis_client *RedisClient) GetSet(key string, value any) string {
	result, err := redis_client.pool.GetSet(key, value).Result()
	if err != nil {
		log.Printf("Error when getting set key %s with value %s: %v", key, value, err)
		return Empty
	}
	return result
}

//SetNx sets the value of a key, only if the key does not exist.
func (redis_client *RedisClient) SetNx(key string, value string, expiration time.Duration) bool {
	result, err := redis_client.pool.SetNX(key, value, expiration).Result()
	if err != nil {
		log.Printf("Error when setting key if not exist %s: %v", key, err)
		return false
	}
	return result
}

// SetEx sets the value of a key, with an expiration time, only if the key exists.
func (redis_client *RedisClient) SetEx(key string, value string, expiration time.Duration) bool {
	result, err := redis_client.pool.SetXX(key, value, expiration).Result()
	if err != nil {
		log.Printf("Error when setting key if exist %s: %v", key, err)
		return false
	}
	return result
}

// MGet returns the values of multiple keys.
func (redis_client *RedisClient) MGet(keys ...string) []any {
	result, err := redis_client.pool.MGet(keys...).Result()
	if err != nil {
		log.Printf("Error when getting multiple keys %s: %v", fmt.Sprint(keys), err)
		return nil
	}
	return result
}

// MSet sets multiple keys to multiple values.
func (redis_client *RedisClient) MSet(pairs ...string) string {
	result, err := redis_client.pool.MSet(pairs).Result()
	if err != nil {
		log.Printf("Error when setting multiple keys %s: %v", fmt.Sprint(pairs), err)
		return Empty
	}
	return result
}

// IncrBy increments the value of a key by the given amount.
func (redis_client *RedisClient) IncrBy(key string, value int) int64 {
	result, err := redis_client.pool.IncrBy(key, int64(value)).Result()
	if err != nil {
		log.Printf("Error when incrementing key %s by %d: %v", key, value, err)
		return 0
	}
	return result
}

// Incr increments the value of a key by 1.
func (redis_client *RedisClient) Incr(key string) int64 {
	result, err := redis_client.pool.Incr(key).Result()
	if err != nil {
		log.Printf("Error when incrementing key %s by 1: %v", key, err)
		return 0
	}
	return result
}

// DecrBy decrements the value of a key by the given amount.
func (redis_client *RedisClient) DecrBy(key string, value int) int64 {
	result, err := redis_client.pool.DecrBy(key, int64(value)).Result()
	if err != nil {
		log.Printf("Error when decrementing key %s by %d: %v", key, value, err)
		return 0
	}
	return result
}

// Decr decrements the value of a key by 1.
func (redis_client *RedisClient) Decr(key string) int64 {
	result, err := redis_client.pool.Decr(key).Result()
	if err != nil {
		log.Printf("Error when decrementing key %s by 1: %v", key, err)
		return 0
	}
	return result
}

// Del deletes multiple keys.
func (redis_client *RedisClient) Del(keys ...string) int64 {
	result, err := redis_client.pool.Del(keys...).Result()
	if err != nil {
		log.Printf("Error when deleting multiple keys %s: %v", fmt.Sprint(keys), err)
		return FAIL
	}
	return result
}

// Expires sets an expiration on a key.
func (redis_client *RedisClient) Expires(key string, expiration time.Duration) bool {
	result, err := redis_client.pool.Expire(key, expiration).Result()
	if err != nil {
		log.Printf("Error when setting expiration on key %s: %v", key, err)
		return false
	}
	return result
}
