package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Student struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name"`
	Email       string               `json:"email" bson:"email"`
	Password    string               `json:"password,omitempty" bson:"password"`
	Courses     []primitive.ObjectID `json:"courses" bson:"courses"`
	Attendances []primitive.ObjectID `json:"attendances" bson:"attendances"`
}
