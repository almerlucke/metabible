package store

import (
	"context"
	"errors"
	"github.com/almerlucke/metabible/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	ConnectionString string `json:"connectionString"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Database         string `json:"database"`
}

type Store struct {
	client *mongo.Client
	cfg    *Config
}

func New(cfg *Config) (*Store, error) {
	var store Store

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.ConnectionString).SetAuth(options.Credential{
		Username: cfg.Username,
		Password: cfg.Password,
	}))

	if err != nil {
		return nil, err
	}

	store.client = client
	store.cfg = cfg

	if err = store.createIndexes(ctx); err != nil {
		return nil, err
	}

	return &store, nil
}

func (store *Store) DB() *mongo.Database {
	return store.client.Database(store.cfg.Database)
}

func (store *Store) Collection(col string) *mongo.Collection {
	return store.DB().Collection(col)
}

func (store *Store) Close() error {
	return store.client.Disconnect(context.Background())
}

func (store *Store) AddDocuments(ctx context.Context, colName string, docs []any) error {
	col := store.DB().Collection(colName)

	cnt, err := col.CountDocuments(ctx, bson.M{}, nil)
	if err != nil {
		return err
	}

	if cnt > 0 {
		return nil
	}

	_, err = col.InsertMany(ctx, docs, nil)
	if err != nil {
		return err
	}

	return nil
}

func Seed[T any](store *Store, ctx context.Context, file string, collection string) error {
	input, err := models.LoadDocs[T](file)
	if err != nil {
		return err
	}

	docs := make([]any, len(input))
	for i, d := range input {
		docs[i] = d
	}

	return store.AddDocuments(ctx, collection, docs)
}

func (store *Store) GetAllPlaces(ctx context.Context) ([]*models.Place, error) {
	var result []*models.Place

	cursor, err := store.Collection("places").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &result)

	return result, err
}

func (store *Store) GetPlaceByLookup(ctx context.Context, placeLookup string) (*models.Place, error) {
	var result models.Place

	err := store.Collection("places").FindOne(ctx, bson.M{"place_lookup": placeLookup}, nil).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (store *Store) GetPlaceByID(ctx context.Context, placeID string) (*models.Place, error) {
	var result models.Place

	err := store.Collection("places").FindOne(ctx, bson.M{"_id": placeID}, nil).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (store *Store) GetVerseByID(ctx context.Context, id string) (*models.Verse, error) {
	var result models.Verse

	err := store.Collection("verses").FindOne(ctx, bson.M{"_id": id}, nil).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (store *Store) GetPersonByLookup(ctx context.Context, personLookup string) (*models.Person, error) {
	var result models.Person

	err := store.Collection("people").FindOne(ctx, bson.M{"person_lookup": personLookup}, nil).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}
