package main

import (
	"context"
	"github.com/allegro/bigcache/v2"
	"github.com/gin-gonic/gin"
	"idempotency/internal/api"
	"idempotency/internal/api/handler"
	"idempotency/internal/api/middleware"
	"idempotency/internal/cache"
	"idempotency/internal/infra"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	db, err := infra.OpenSqliteConnection("user.db")
	if err != nil {
		log.Fatalln(db)
	}

	defaultTTL := time.Minute * 10

	r := infra.NewGormUserRepository(db)
	h := handler.NewUserHandler(r)
	p := infra.NewBigCacheProvider(bigcache.DefaultConfig(defaultTTL))
	c := cache.NewIdempotencyCache(p, defaultTTL)

	engine := gin.Default()
	engine.Use(middleware.Tracing, middleware.Idempotency(c, middleware.DefaultIntent))
	engine.GET("/user/:id", h.GetUser)
	engine.POST("/user", h.CreateUser)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-shutdown
		log.Println("shutting down...")
		cancel()
	}()

	if err := api.StartServer(ctx, engine, ":8080"); err != nil {
		log.Fatalln(err)
	}
}
