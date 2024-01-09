package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	contactCollection = "contacts"
	database = "contact_db"
)

type ContactRepository struct {
	client *mongo.Client
}

func NewContactRepository(client *mongo.Client) *ContactRepository {
	return &ContactRepository{client}
}

func (r *ContactRepository) Save(ctx context.Context, contact *Contact) error {
	contact.ID = primitive.NewObjectID().Hex()
	_, err := r.client.Database(database).Collection(contactCollection).InsertOne(ctx, contact)
	return err
}

func (r *ContactRepository) FindByID(ctx context.Context, id string) (*Contact, error) {
	var contact Contact
	err := r.client.Database(database).Collection(contactCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&contact)
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

