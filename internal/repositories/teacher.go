package repositories

import (
	"context"
	"fmt"
	"money-minder/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeacherRepo struct {
	MongoCollection *mongo.Collection
}

func (r *TeacherRepo) InsertTeacher(teacher *types.Teacher) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), teacher)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *TeacherRepo) DeleteTeacher(teacherID string) error {
	id, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		return err
	}
	_, err = r.MongoCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *TeacherRepo) FindTeacherByID(teacherID string) (*types.Teacher, error) {
	id, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		return nil, err
	}
	var teacher types.Teacher
	err = r.MongoCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&teacher)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &teacher, nil
}

func (r *TeacherRepo) FindAllTeachers() ([]types.Teacher, error) {
	cursor, err := r.MongoCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var teachers []types.Teacher
	if err = cursor.All(context.Background(), &teachers); err != nil {
		return nil, err
	}
	return teachers, nil
}

func (r *TeacherRepo) AddCourse(teacherID string, course types.Course) error {
	teacherObjID, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		return fmt.Errorf("invalid teacher ID: %w", err)
	}
	_, err = r.MongoCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": teacherObjID},
		bson.M{"$addToSet": bson.M{"courses": course}},
	)
	return err
}

func (r *TeacherRepo) RemoveCourse(teacherID string, courseID string) error {
	teacherObjID, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		return fmt.Errorf("invalid teacher ID: %w", err)
	}
	courseObjID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return fmt.Errorf("invalid course ID: %w", err)
	}
	_, err = r.MongoCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": teacherObjID},
		bson.M{"$pull": bson.M{"courses": bson.M{"_id": courseObjID}}},
	)
	return err
}

func (r *TeacherRepo) GetAllCoursesByTeacherID(teacherID string) ([]types.Course, error) {
	id, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		return nil, err
	}
	var teacher types.Teacher
	err = r.MongoCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&teacher)
	if err != nil {
		return nil, err
	}
	return teacher.Courses, nil
}

func (r *TeacherRepo) FindTeacherByEmail(email string) (*types.Teacher, error) {
	var teacher types.Teacher
	err := r.MongoCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&teacher)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &teacher, nil
}
