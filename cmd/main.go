package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/crisywini/pinta-api/config"
	"github.com/crisywini/pinta-api/internal/bootstrap"
	"github.com/crisywini/pinta-api/internal/handler"
	"github.com/crisywini/pinta-api/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	repos, err := bootstrap.NewRepositories(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := repos.Close(closeCtx); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	itemService := service.NewItemService(repos.Items)
	outfitService := service.NewOutfitService(repos.Outfits, repos.Items)

	itemHandler := handler.NewItemHandler(itemService)
	outfitHandler := handler.NewOutfitHandler(outfitService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/items", itemHandler.PostItem)
	r.GET("/items/:id", itemHandler.GetItemByID)
	r.GET("/items", itemHandler.GetAllItems)
	r.PUT("/items/:id", itemHandler.PutItem)
	r.DELETE("/items/:id", itemHandler.DeleteItem)

	r.POST("/outfits", outfitHandler.PostOutfit)
	r.GET("/outfits/:id", outfitHandler.GetOutfitById)
	r.GET("/outfits", outfitHandler.GetAllOutfits)
	r.PUT("/outfits/:id", outfitHandler.PutOutfit)
	r.DELETE("/outfits/:id", outfitHandler.DeleteOutfitById)

	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done() // espera Ctrl+C / SIGTERM
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
