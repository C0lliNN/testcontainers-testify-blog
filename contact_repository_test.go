package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ContactRepositorySuite struct {
	RepositorySuite
	repo *ContactRepository
}

func (s *ContactRepositorySuite) SetupSuite() {
	s.RepositorySuite.SetupSuite()

	s.repo = NewContactRepository(s.client)
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
