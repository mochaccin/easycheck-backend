package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Code     string             `json:"code" bson:"code"`
	Teacher  primitive.ObjectID `json:"teacher" bson:"teacher,omitempty"`
	Students []Student          `json:"students,omitempty" bson:"students,omitempty"`
}
