package repository

type Contact struct {
	ID string `bson:"_id"`
	Name string
	Phone string
	Email string
}