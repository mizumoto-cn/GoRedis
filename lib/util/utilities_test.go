package util

import (
	"log"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	pool = NewRedisClient("localhost:6379", 0)
)

// TestRedisClient tests the RedisClient.
func TestRedisClient(t *testing.T) {

	str := "Mizumoto"
	set := pool.Set("test", str, 0)
	log.Println(set)
	get := pool.Get("test")
	assert.Equal(t, str, get)

	getset := pool.GetSet("test", str+"2")
	log.Println(getset)
	pool.Set("neumi", "nemuinya", 0)
	mget := pool.MGet("test", "nemui")
	for k, v := range mget {
		log.Printf("%d: %s", k, v)
	}
	pool.Expires("test", 300)

	keys := pool.Keys("miz")
	for i, key := range keys {
		log.Printf("index : %d, key : %s", i, key)
	}
	key := str + strconv.Itoa(int(rand.Int63()))
	pool.SetNx(key, "3432", 100)
	log.Println(pool.Get(key))
	log.Println("22222")
	i := pool.MGet(str, key)
	for _, i3 := range i {
		log.Println(i3)
	}
	pool.MSet("aa", "111", "aabbbb", "22222")
	pool.Incr("sum")
	pool.IncrBy("sum", 10)
	log.Println(pool.Get("sum"))
	strings := pool.Keys("aa*")
	for _, s := range strings {
		log.Println(s)
	}
	pool.Decr("aa")
	pool.Expires("sum", 100)
}
