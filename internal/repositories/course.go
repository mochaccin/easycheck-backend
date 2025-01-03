package repositories

import (
	"context"
	"fmt"
	"money-minder/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseRepo struct {
	MongoCollection *mongo.Collection
}

func (r *CourseRepo) InsertCourse(course *types.Course) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), course)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *CourseRepo) DeleteCourse(courseID string) error {
	id, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return err
	}
	_, err = r.MongoCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *CourseRepo) FindCourseByID(courseID string) (*types.Course, error) {
	id, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return nil, err
	}
	var course types.Course
	err = r.MongoCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&course)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (r *CourseRepo) ChangeTeacher(courseID string, teacherID string) error {
	courseObjID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return fmt.Errorf("invalid course ID: %w", err)
	}
	teacherObjID, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		return fmt.Errorf("invalid teacher ID: %w", err)
	}
	_, err = r.MongoCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": courseObjID},
		bson.M{"$set": bson.M{"teacher": teacherObjID}},
	)
	return err
}

func (r *CourseRepo) AddStudent(courseID string, student types.Student) error {
	courseObjID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return fmt.Errorf("invalid course ID: %w", err)
	}
	_, err = r.MongoCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": courseObjID},
		bson.M{"$addToSet": bson.M{"students": student}},
	)
	return err
}

func (r *CourseRepo) RemoveStudent(courseID string, studentID string) error {
	courseObjID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return fmt.Errorf("invalid course ID: %w", err)
	}
	studentObjID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return fmt.Errorf("invalid student ID: %w", err)
	}
	_, err = r.MongoCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": courseObjID},
		bson.M{"$pull": bson.M{"students": bson.M{"_id": studentObjID}}},
	)
	return err
}

func (r *CourseRepo) GetAllStudentsByCourseID(courseID string) ([]types.Student, error) {
	id, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return nil, err
	}
	var course types.Course
	err = r.MongoCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&course)
	if err != nil {
		return nil, err
	}
	return course.Students, nil
}

func (r *CourseRepo) FindAllCourses() ([]types.Course, error) {
	cursor, err := r.MongoCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var courses []types.Course
	if err = cursor.All(context.Background(), &courses); err != nil {
		return nil, err
	}
	return courses, nil
}
