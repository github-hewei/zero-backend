package mongodb

import (
	"context"
	"time"
	"zero-backend/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Conn struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewConn(cfg *config.Config) *Conn {
	if !cfg.MongoDB.Enabled {
		return &Conn{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	return &Conn{
		Client: client,
		DB:     client.Database(cfg.MongoDB.Database),
	}
}
