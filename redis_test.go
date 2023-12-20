package belajar_golang_redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   0,
})

func TestConnection(t *testing.T) {
	assert.NotNil(t, client)

	// err := client.Close()
	// assert.Nil(t, err)
}

var ctx = context.Background()

func TestPing(t *testing.T) {
	result, err := client.Ping(ctx).Result()
	assert.Nil(t, err)
	assert.Equal(t, "PONG", result)
}

func TestString(t *testing.T) {
	client.SetEx(ctx, "name", "Aldi Syah Putra", 3*time.Second)

	result, err := client.Get(ctx, "name").Result()
	assert.Nil(t, err)
	assert.Equal(t, "Aldi Syah Putra", result)

	time.Sleep(5 * time.Second)

	result, err = client.Get(ctx, "name").Result()
	assert.NotNil(t, err)
}

func TestList(t *testing.T) {
	client.RPush(ctx, "names", "Aldi")
	client.RPush(ctx, "names", "Syah")
	client.RPush(ctx, "names", "Putra")

	assert.Equal(t, "Aldi", client.LPop(ctx, "names").Val())
	assert.Equal(t, "Syah", client.LPop(ctx, "names").Val())
	assert.Equal(t, "Putra", client.LPop(ctx, "names").Val())

	client.Del(ctx, "names")
}

func TestSet(t *testing.T) {
	client.SAdd(ctx, "students", "Aldi")
	client.SAdd(ctx, "students", "Aldi")
	client.SAdd(ctx, "students", "Syah")
	client.SAdd(ctx, "students", "Syah")
	client.SAdd(ctx, "students", "Putra")
	client.SAdd(ctx, "students", "Putra")

	assert.Equal(t, int64(3), client.SCard(ctx, "students").Val())
	assert.Equal(t, []string{"Aldi", "Syah", "Putra"}, client.SMembers(ctx, "students").Val())
}

func TestSortedSet(t *testing.T) {
	client.ZAdd(ctx, "scores", redis.Z{Score: 100, Member: "Aldi"})
	client.ZAdd(ctx, "scores", redis.Z{Score: 85, Member: "Budi"})
	client.ZAdd(ctx, "scores", redis.Z{Score: 95, Member: "Joko"})

	assert.Equal(t, []string{"Budi", "Joko", "Aldi"}, client.ZRange(ctx, "scores", 0, -1).Val())
	assert.Equal(t, "Aldi", client.ZPopMax(ctx, "scores").Val()[0].Member)
	assert.Equal(t, "Joko", client.ZPopMax(ctx, "scores").Val()[0].Member)
	assert.Equal(t, "Budi", client.ZPopMax(ctx, "scores").Val()[0].Member)
}

func TestHash(t *testing.T) {
	client.HSet(ctx, "user:1", "id", "1")
	client.HSet(ctx, "user:1", "name", "Aldi")
	client.HSet(ctx, "user:1", "email", "aldi@example.com")

	user := client.HGetAll(ctx, "user:1").Val()
	assert.Equal(t, "1", user["id"])
	assert.Equal(t, "Aldi", user["name"])
	assert.Equal(t, "aldi@example.com", user["email"])

	client.Del(ctx, "user:1")
}
