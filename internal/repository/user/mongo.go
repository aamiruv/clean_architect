package user

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/amirzayi/clean_architect/internal/domain"
)

const userCollectionName = "user"

type userMongoRepo struct {
	db *mongo.Collection
}

func NewUserMongoRepository(db *mongo.Database) *userMongoRepo {
	return &userMongoRepo{db: db.Collection(userCollectionName)}
}

func (r *userMongoRepo) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.InsertOne(ctx, user)
	return err
}
