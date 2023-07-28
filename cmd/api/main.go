package main

import (
	"context"
	journeyInternal "github.com/dportaluppi/journey-api/internal/journey"
	"github.com/dportaluppi/journey-api/pkg/journey"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	// Set up MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the primary
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	// Initialize the journey repository, services, and handler
	repo := journeyInternal.NewMongoRepo(client)
	handler := journeyInternal.NewHandler(
		journey.NewGetter(repo),
		journey.NewCreator(repo),
		journey.NewUpdater(repo),
		journey.NewDeleter(repo),
	)

	// Set up the HTTP server
	r := gin.Default()
	api := r.Group("/journeys")
	api.GET("", handler.GetJourneys)
	{
		api.GET("/:id", handler.GetJourneyByID)
		api.POST("", handler.CreateJourney)
		api.PUT("/:id", handler.UpdateJourney)
		api.DELETE("/:id", handler.DeleteJourney)
	}

	// Start the HTTP server
	if err := r.Run(":8060"); err != nil {
		log.Fatal(err)
	}
}
