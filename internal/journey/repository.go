package journey

import (
	"context"
	"fmt"
	"github.com/dportaluppi/journey-api/pkg/journey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type MongoRepo struct {
	client         *mongo.Client
	databaseName   string
	collectionName string
}

func NewMongoRepo(client *mongo.Client) *MongoRepo {
	return &MongoRepo{
		client:         client,
		databaseName:   "journey_api",
		collectionName: "account_journeys",
	}
}

func (r *MongoRepo) GetJourneys(ctx context.Context, filter *journey.Filter, sortBy string) ([]journey.Journey, error) {
	var journeys []journey.Journey

	collection := r.client.Database(r.databaseName).Collection(r.collectionName)

	findFilter := bson.M{}

	if filter.Storefront != "" {
		findFilter["storefront"] = filter.Storefront
	}

	if !filter.Date.IsZero() {
		findFilter["startAt"] = bson.M{"$lte": filter.Date}
		findFilter["$or"] = bson.A{
			bson.M{"endAt": bson.M{"$gte": filter.Date}},
			bson.M{"endAt": bson.M{"$eq": nil}},
		}
	}

	if len(filter.Audiences) > 0 {
		findFilter["audiences"] = bson.M{"$in": filter.Audiences}
	}

	if len(filter.Channels) > 0 {
		findFilter["channels.id"] = bson.M{"$in": filter.Channels}
	}

	findOptions := options.Find()

	if sortBy != "" {
		sortDirection := 1
		if strings.HasPrefix(sortBy, "-") {
			sortBy = strings.TrimPrefix(sortBy, "-")
			sortDirection = -1
		}
		findOptions.SetSort(bson.D{{sortBy, sortDirection}})
	}

	cur, err := collection.Find(ctx, findFilter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var j journey.Journey
		err := cur.Decode(&j)
		if err != nil {
			return nil, err
		}
		journeys = append(journeys, j)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return journeys, nil
}

func (r *MongoRepo) GetByID(ctx context.Context, id string) (*journey.Journey, error) {
	var j journey.Journey
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	collection := r.client.Database(r.databaseName).Collection(r.collectionName)
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&j)
	return &j, err
}

func (r *MongoRepo) Create(ctx context.Context, j *journey.Journey) (*journey.Journey, error) {
	collection := r.client.Database(r.databaseName).Collection(r.collectionName)
	res, err := collection.InsertOne(ctx, j)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}

	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("inserted id is not of type primitive.ObjectID")
	}

	j.ID = objectID

	return j, nil
}

func (r *MongoRepo) Update(ctx context.Context, id string, j *journey.Journey) (*journey.Journey, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	collection := r.client.Database(r.databaseName).Collection(r.collectionName)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": j})
	if err != nil {
		return nil, err
	}

	var updatedJourney journey.Journey
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedJourney)
	if err != nil {
		return nil, err
	}

	return &updatedJourney, nil
}

func (r *MongoRepo) Delete(ctx context.Context, id string) error {
	collection := r.client.Database(r.databaseName).Collection(r.collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
