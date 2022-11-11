package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func Test_redis(t *testing.T) {
	u := &User{}
	u.red = NewRedisCache("postgres:6379", "")
	c := context.Background()
	value, err := u.red.Get(c, "demo_key").Result()
	if err == redis.Nil {
		_, err := u.red.Set(c, "demo_key", 1, 10*time.Second).Result()
		if err != nil {
			log.Panicln(err)
			return
		}
		value, _ = u.red.Get(c, "demo_key").Result()
	}
	log.Println("demo_key: ", value)
}

func Test_ScanToCache(t *testing.T) {
	data := make(chan string, 100)
	wg := &sync.WaitGroup{}
	cCtx, cancel := context.WithCancel(context.Background())
	file, err := os.Open("name.txt")
	if err != nil {
		log.Panicln(err)
	}
	go func() {
		for {
			select {
			case <-cCtx.Done():
				return
			case name := <-data:
				log.Println("name: ", name)
				wg.Done()
				continue
			}
		}
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data <- scanner.Text()
		wg.Add(1)
	}
	wg.Wait()
	file.Close()
	cancel()
	close(data)
}
