package main

import (
	"context"
	"log"
	"os"

	"github.com/crisywini/pinta-api/internal/handler"
	"github.com/crisywini/pinta-api/internal/repository"
	"github.com/crisywini/pinta-api/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	var databaseUri string
	if databaseUri = os.Getenv("MONGODB_URI"); databaseUri == "" {
		log.Fatal("You must set the MONGODB_URI for mongodb connection")
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(databaseUri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	database := client.Database("pintapp")
	itemRepository := repository.NewItemRepository(database)
	itemService := service.NewItemService(itemRepository)
	outfitRepository := repository.NewOutfitRepository(database)
	outfitService := service.NewOutfitService(outfitRepository, itemRepository)

	itemHandler := handler.NewItemHandler(itemService)
	outfitHandler := handler.NewOutfitHandler(outfitService)

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

	r.Run(":8080")
}
