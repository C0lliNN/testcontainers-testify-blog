package repository

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepositorySuite struct {
	suite.Suite
	client *mongo.Client
	container *mongodb.MongoDBContainer
}

func (s *RepositorySuite) SetupSuite() {
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
}

func (s *RepositorySuite) TearDownSuite() {
	ctx := context.Background()
	
	err := s.client.Disconnect(ctx)
	s.Require().NoError(err)

	err = s.container.Terminate(ctx)
	s.Require().NoError(err)
}
