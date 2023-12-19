package mongodb

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rasulov-emirlan/poc-auth/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepoCombiner struct {
	conn *mongo.Client
}

func NewRepoCombiner(ctx context.Context, cfg config.Config, logger *slog.Logger) (RepoCombiner, error) {
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.MongoDB.GetURI()))
	if err != nil {
		return RepoCombiner{}, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	if err := conn.Ping(ctx, nil); err != nil {
		return RepoCombiner{}, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return RepoCombiner{
		conn: conn,
	}, nil
}

func (r RepoCombiner) Close(ctx context.Context) error {
	return r.conn.Disconnect(ctx)
}

func (r RepoCombiner) Check(ctx context.Context) error {
	return r.conn.Ping(ctx, nil)
}
