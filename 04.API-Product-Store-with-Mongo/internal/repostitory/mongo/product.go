package mongo

import (
	"04.API-Product-Store-with-Mongo/internal/core"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) *ProductRepository {
	return &ProductRepository{collection: collection}
}

func (repository *ProductRepository) GetAll(ctx context.Context) ([]*core.Product, error) {
	cursor, err := repository.collection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	products := make([]*core.Product, 0)

	err = cursor.All(ctx, &products)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (repository *ProductRepository) GetById(ctx context.Context, id string) (*core.Product, error) {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	product := &core.Product{}

	err = repository.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(product)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (repository *ProductRepository) Save(ctx context.Context, product *core.Product) (*core.Product, error) {
	result, err := repository.collection.InsertOne(ctx, product)

	if err != nil {
		return nil, err
	}

	product.ID = result.InsertedID.(primitive.ObjectID)

	return product, nil
}
