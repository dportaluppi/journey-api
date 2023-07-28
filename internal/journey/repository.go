package journey

import (
	"context"
	"fmt"
	"github.com/dportaluppi/journey-api/pkg/journey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	client *mongo.Client
}

func NewMongoRepo(client *mongo.Client) *MongoRepo {
	return &MongoRepo{client: client}
}

func (r *MongoRepo) GetJourneys(ctx context.Context, filter *journey.Filter, sortBy string) ([]journey.Journey, error) {
	var journeys []journey.Journey

	collection := r.client.Database("mydatabase").Collection("journeys")

	findFilter := bson.M{}

	if filter.Storefront != "" {
		findFilter["storefront"] = filter.Storefront
	}

	if !filter.Date.IsZero() {
		findFilter["startAt"] = bson.M{"$lte": filter.Date}
		findFilter["endAt"] = bson.M{"$gte": filter.Date}
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
		var journey journey.Journey
		err := cur.Decode(&journey)
		if err != nil {
			return nil, err
		}
		journeys = append(journeys, journey)
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
	collection := r.client.Database("mydatabase").Collection("journeys")
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&j)
	return &j, err
}

func (r *MongoRepo) Create(ctx context.Context, j *journey.Journey) error {
	collection := r.client.Database("mydatabase").Collection("journeys")
	res, err := collection.InsertOne(ctx, j)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("inserted id is not of type primitive.ObjectID")
	}

	j.ID = objectID
	return nil
}

func (r *MongoRepo) Update(ctx context.Context, id string, j *journey.Journey) error {
	collection := r.client.Database("mydatabase").Collection("journeys")
	_, err := collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": j})
	return err
}

func (r *MongoRepo) Delete(ctx context.Context, id string) error {
	collection := r.client.Database("mydatabase").Collection("journeys")
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func addFilterCondition(filter bson.M, condition string, value interface{}) {
	if value != nil && !isEmpty(value) {
		filter[condition] = value
	}
}

func isEmpty(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return v == ""
	case time.Time:
		return v.IsZero()
	case []string:
		return len(v) == 0
	default:
		return false
	}
}
