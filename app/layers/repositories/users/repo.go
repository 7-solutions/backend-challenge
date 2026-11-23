package users

import (
	"context"
	"time"

	"github.com/chaihanij/backend-challenge/app/database"
	"github.com/chaihanij/backend-challenge/app/entities"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	CollectionUser = "users"
)

type Repo interface {
	CountUser(ctx context.Context) (int64, error)
	CreateUser(ctx context.Context, input *entities.User) (*entities.User, error)
	DeleteUser(ctx context.Context, id string) error
	FindOneUser(ctx context.Context, id string) (*entities.User, error)
	FindUser(ctx context.Context, input *entities.UserFilter) ([]entities.User, error)
	UpdateUser(ctx context.Context, input *entities.User) error
}

type repo struct {
	col     database.MongoCollection
	timeout time.Duration
}

func New(mongoClient *mongo.Client, mongoDBName string, timeout time.Duration) Repo {
	col := mongoClient.Database(mongoDBName).Collection(CollectionUser)
	return &repo{
		col:     col,
		timeout: timeout,
	}
}

func Wrap(col database.MongoCollection, timeout time.Duration) Repo {
	return &repo{
		col:     col,
		timeout: timeout,
	}
}
