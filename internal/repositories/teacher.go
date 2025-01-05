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

func (r *TeacherRepo) InsertTeacher(usr *types.Teacher) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne((context.Background()), usr)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TeacherRepo) DeleteTeacher(usr *types.Teacher) (interface{}, error) {
	result, err := r.MongoCollection.DeleteOne((context.Background()), usr)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TeacherRepo) FindTeacherByID(usrID string) (*types.Teacher, error) {
	id, err := primitive.ObjectIDFromHex(usrID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}

	var usr types.Teacher

	err = r.MongoCollection.FindOne(context.Background(), filter).Decode(&usr)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &usr, nil
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

func (r *TeacherRepo) FindAllTeachers() ([]types.Teacher, error) {
	results, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var usrs []types.Teacher

	err = results.All(context.Background(), &usrs)
	if err != nil {
		return nil, fmt.Errorf("Find all uses results decode error %s", err.Error())
	}

	return usrs, nil
}

func (r *TeacherRepo) UpdateName(usrID string, newName string) error {
	id, err := primitive.ObjectIDFromHex(usrID)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"name", newName}}}}

	_, err = r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *TeacherRepo) AddCourse(TeacherID string, CourseID string, CourseRepo *CourseRepo) error {

	TeacherObjectID, err := primitive.ObjectIDFromHex(TeacherID)
	if err != nil {
		return fmt.Errorf("invalid Teacher ID: %w", err)
	}

	Course, err := CourseRepo.FindCourseByID(CourseID)
	if err != nil {
		return fmt.Errorf("failed to find Course: %w", err)
	}
	if Course == nil {
		return fmt.Errorf("Course not found")
	}

	filter := bson.D{{"_id", TeacherObjectID}}
	update := bson.D{{"$push", bson.D{{"courses", Course}}}}

	_, err = r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to add Course to Teacher: %w", err)
	}

	return nil
}

func (r *TeacherRepo) RemoveCourse(TeacherID string, CourseID string, CourseRepo *CourseRepo) error {

	TeacherObjectID, err := primitive.ObjectIDFromHex(TeacherID)
	if err != nil {
		return fmt.Errorf("invalid Teacher ID: %w", err)
	}

	Course, err := CourseRepo.FindCourseByID(CourseID)
	if err != nil {
		return fmt.Errorf("failed to find Course: %w", err)
	}
	if Course == nil {
		return fmt.Errorf("Course not found")
	}

	filter := bson.D{{"_id", TeacherObjectID}}
	update := bson.D{{"$pull", bson.D{{"courses", bson.D{{"$eq", Course}}}}}}

	_, err = r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete Course from Teacher: %w", err)
	}

	return nil
}
