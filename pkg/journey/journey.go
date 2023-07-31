package journey

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Journey struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name            string             `bson:"name" json:"name"`
	Storefront      string             `bson:"storefront" json:"storefront"`
	Priority        int                `bson:"priority" json:"priority"`
	Audiences       []string           `bson:"audiences" json:"audiences"`
	Recommendations []Recommendation   `bson:"recommendations" json:"recommendations"`
	Channels        []Channel          `bson:"channels" json:"channels"`
	Frequency       Frequency          `bson:"frequency" json:"frequency"`
	StartAt         *time.Time         `bson:"startAt" json:"startAt"`
	EndAt           *time.Time         `bson:"endAt" json:"endAt"`
}

type Recommendation struct {
	Type     string   `bson:"type" json:"type"`
	Products []string `bson:"products" json:"products"`
}

type Channel struct {
	ID      string   `bson:"id" json:"id"`
	Options []Option `bson:"options" json:"options"`
}

type Option struct {
	Type  string `bson:"type" json:"type"`
	Value string `bson:"value" json:"value"`
}

type Frequency struct {
	Cron     string `bson:"cron" json:"cron"`
	Timezone string `bson:"timezone" json:"timezone"`
}

type Filter struct {
	Storefront string
	Audiences  []string
	Channels   []string
	Date       time.Time
}

type Getter interface {
	GetJourneysByCriteria(ctx context.Context, storefront string, audiences, channels []string, date time.Time, sort string) ([]Journey, error)
	GetJourneyByID(ctx context.Context, id string) (*Journey, error)
}

type Creator interface {
	CreateJourney(ctx context.Context, j *Journey) (*Journey, error)
}

type Updater interface {
	UpdateJourney(ctx context.Context, id string, j *Journey) (*Journey, error)
}

type Deleter interface {
	DeleteJourney(ctx context.Context, id string) error
}

type Repository interface {
	GetJourneys(ctx context.Context, filter *Filter, sortBy string) ([]Journey, error)
	GetByID(ctx context.Context, id string) (*Journey, error)
	Create(ctx context.Context, j *Journey) (*Journey, error)
	Update(ctx context.Context, id string, j *Journey) (*Journey, error)
	Delete(ctx context.Context, id string) error
}
