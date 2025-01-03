package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	GetCollection(name string) *mongo.Collection
}

type service struct {
	db *mongo.Client
}

var (
	instance *service
	once     sync.Once
	uri      = os.Getenv("DB_URI")
)

func New() Service {
	once.Do(func() {
		if uri == "" {
			log.Fatal("DB_URI is not set in the environment")
		}

		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err != nil {
			log.Fatal(err)
		}

		instance = &service{
			db: client,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := instance.db.Ping(ctx, nil); err != nil {
			log.Fatalf(fmt.Sprintf("db connection error: %v", err))
		}
	})
	return instance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) GetCollection(name string) *mongo.Collection {
	return s.db.Database("easycheck").Collection(name)
}
