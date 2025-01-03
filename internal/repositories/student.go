package repositories

import (
	"context"
	"fmt"
	"money-minder/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudentRepo struct {
	MongoCollection      *mongo.Collection
	CourseCollection     *mongo.Collection
	AttendanceCollection *mongo.Collection
}

func (r *StudentRepo) InsertStudent(student *types.Student) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), student)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *StudentRepo) DeleteStudent(studentID string) error {
	id, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return err
	}
	_, err = r.MongoCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *StudentRepo) FindStudentByID(studentID string) (*types.Student, error) {
	id, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, err
	}
	var student types.Student
	err = r.MongoCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepo) FindAllStudents() ([]types.Student, error) {
	cursor, err := r.MongoCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var students []types.Student
	if err = cursor.All(context.Background(), &students); err != nil {
		return nil, err
	}
	return students, nil
}

func (r *StudentRepo) AddCourse(studentID string, courseID string) error {
	studentObjID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return fmt.Errorf("invalid student ID: %w", err)
	}

	courseObjID, err := primitive.ObjectIDFromHex(courseID)

	if err != nil {
		return fmt.Errorf("invalid course ID: %w", err)
	}
	_, err = r.MongoCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": studentObjID},
		bson.M{"$addToSet": bson.M{"courses": courseObjID}},
	)
	return err
}

func (r *StudentRepo) RemoveCourse(studentID string, courseID string) error {
	studentObjID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return fmt.Errorf("invalid student ID: %w", err)
	}
	courseObjID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return fmt.Errorf("invalid course ID: %w", err)
	}
	_, err = r.MongoCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": studentObjID},
		bson.M{"$pull": bson.M{"courses": bson.M{"_id": courseObjID}}},
	)
	return err
}

func (r *StudentRepo) GetAllCoursesByStudentID(studentID string) ([]types.Course, error) {
	id, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, err
	}
	var student types.Student
	err = r.MongoCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&student)
	if err != nil {
		return nil, err
	}

	var courses []types.Course
	cursor, err := r.CourseCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": student.Courses}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &courses); err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *StudentRepo) GetAllAttendancesByStudentID(studentID string) ([]types.Attendance, error) {
	id, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, err
	}
	var student types.Student
	err = r.MongoCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&student)
	if err != nil {
		return nil, err
	}

	var attendances []types.Attendance
	cursor, err := r.AttendanceCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": student.Attendances}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &attendances); err != nil {
		return nil, err
	}

	return attendances, nil
}

func (r *StudentRepo) FindStudentByEmail(email string) (*types.Student, error) {
	var student types.Student
	err := r.MongoCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}
