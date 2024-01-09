package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContactRepositorySuite struct {
	suite.Suite
	client *mongo.Client
	container *mongodb.MongoDBContainer
	repo *ContactRepository
}

func (s *ContactRepositorySuite) SetupSuite() {
	ctx := context.Background()

	var err error
	s.container, err = mongodb.RunContainer(ctx, testcontainers.WithImage("mongo:6"))
	s.Require().NoError(err)

	mongodbPort, err := s.container.MappedPort(ctx, "27017")
	s.Require().NoError(err)

	mongodbHost, err := s.container.Host(ctx)
	s.Require().NoError(err)

	s.client, err = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", mongodbHost, mongodbPort.Port())))
	s.Require().NoError(err)

	err = s.client.Connect(ctx)
	s.Require().NoError(err)

	s.repo = NewContactRepository(s.client)
}

func (s *ContactRepositorySuite) TearDownSuite() {
	ctx := context.Background()
	
	err := s.client.Disconnect(ctx)
	s.Require().NoError(err)

	err = s.container.Terminate(ctx)
	s.Require().NoError(err)
}

func (s *ContactRepositorySuite) TearDownTest() {
	ctx := context.Background()

	err := s.client.Database(database).Drop(ctx)
	s.Require().NoError(err)
}

func (s *ContactRepositorySuite) TestSave() {
	ctx := context.Background()

	contact := &Contact{
		Name: "John Doe",
		Phone: "08123456789",
		Email: "test@test.com",
	}

	err := s.repo.Save(ctx, contact)
	s.Require().NoError(err)

	s.NotEmpty(contact.ID)

	persistedContact, err := s.repo.FindByID(ctx, contact.ID)
	s.Require().NoError(err)
	s.Equal(contact, persistedContact)
}

func TestContactRepositorySuite(t *testing.T) {
	suite.Run(t, new(ContactRepositorySuite))
}
