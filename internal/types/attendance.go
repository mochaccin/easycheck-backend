package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attendance struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CourseID  primitive.ObjectID `json:"courseId" bson:"course_id"`
	StudentID primitive.ObjectID `json:"studentId" bson:"student_id"`
	Date      time.Time          `json:"date" bson:"date"`
	Present   bool               `json:"present" bson:"present"`
}
