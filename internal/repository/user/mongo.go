package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/pkg/paginate"
	"github.com/google/uuid"
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

func (r *userMongoRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user domain.User
	err := r.db.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	return user, err
}

func (r *userMongoRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := r.db.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return user, err
}
func (r *userMongoRepo) List(ctx context.Context, pagination *paginate.Pagination) ([]domain.User, error) {
	var users []domain.User
	cursor, err := r.db.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userMongoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.UpdateOne(ctx,
		bson.M{"id": id},
		bson.M{"$set": bson.M{"status": domain.UserStatusDeleted}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userMongoRepo) Update(ctx context.Context, user domain.User) error {
	res, err := r.db.UpdateOne(ctx,
		bson.M{"id": user.ID},
		bson.M{"$set": bson.M{"name": user.Name, "phone_number": user.PhoneNumber, "email": user.Email, "password": user.Password}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}
