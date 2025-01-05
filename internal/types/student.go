package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Student struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	Courses     []Course           `json:"courses,omitempty" bson:"courses,omitempty"`
	Attendances []Attendance       `json:"attendances,omitempty" bson:"attendances,omitempty"`
}
