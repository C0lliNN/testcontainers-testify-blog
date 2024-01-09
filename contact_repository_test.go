package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func setupTest(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	// Only create a new client if we don't have one yet
	if client == nil {
		mongodbContainer, err := mongodb.RunContainer(ctx, testcontainers.WithImage("mongo:6"))
		if err != nil {
			t.Fatal(err)
		}

		mongodbPort, err := mongodbContainer.MappedPort(ctx, "27017")
		if err != nil {
			t.Fatal(err)
		}

		mongodbHost, err := mongodbContainer.Host(ctx)
		if err != nil {
			t.Fatal(err)
		}

		client, err = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", mongodbHost, mongodbPort.Port())))
		if err != nil {
			t.Fatal(err)
		}

		err = client.Connect(ctx)
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() {
			err = mongodbContainer.Terminate(ctx);
			if err != nil {
				t.Fatal(err)
			}
		})
	}

	t.Cleanup(func() {
		// Delete the database after the test is done to make sure we have a clean state
		err := client.Database(database).Drop(ctx)
		if err != nil {
			t.Fatal(err)
		}
		
	})
}

func TestContactRepository_Save(t *testing.T) {
	setupTest(t)

	ctx := context.Background()
	repo := NewContactRepository(client)

	contact := &Contact{
		Name: "John Doe",
		Phone: "08123456789",
		Email: "test@test.com",
	}

	err := repo.Save(ctx, contact)
	if err != nil {
		t.Fatal(err)
	}

	if contact.ID == "" {
		t.Error("Expected ID to be set")
	}

	persistedContact, err := repo.FindByID(ctx, contact.ID)
	if err != nil {
		t.Fatal(err)
	}

	if persistedContact != contact {
		t.Error("Expected contact to be persisted")
	}
}