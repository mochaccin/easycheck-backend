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

func (r *CourseRepo) InsertCourse(Course *types.Course) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne((context.Background()), Course)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *CourseRepo) DeleteCourse(CourseID string) (*mongo.DeleteResult, error) {
	id, err := primitive.ObjectIDFromHex(CourseID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id}

	result, err := r.MongoCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *CourseRepo) FindCourseByID(CourseID string) (*types.Course, error) {
	id, err := primitive.ObjectIDFromHex(CourseID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}

	var Course types.Course

	err = r.MongoCollection.FindOne(context.Background(), filter).Decode(&Course)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &Course, nil
}

func (r *CourseRepo) FindAllCourses() ([]types.Course, error) {
	results, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var Courses []types.Course

	err = results.All(context.Background(), &Courses)
	if err != nil {
		return nil, fmt.Errorf("Find all uses results decode error %s", err.Error())
	}

	return Courses, nil
}

func (r *CourseRepo) UpdateName(CourseID string, newName string) error {
	id, err := primitive.ObjectIDFromHex(CourseID)
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

func (r *CourseRepo) AddStudent(CourseId string, StudentId string, StudentRepo *StudentRepo) error {

	id, err := primitive.ObjectIDFromHex(CourseId)
	if err != nil {
		return err
	}

	Student, err := StudentRepo.FindStudentByID(StudentId)

	if err != nil {
		return fmt.Errorf("failed to find Course: %w", err)
	}
	if Student == nil {
		return fmt.Errorf("Course not found")
	}

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$push", bson.D{{"students", Student}}}}

	_, err = r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to add Students to Course: %w", err)
	}

	return nil
}

func (r *CourseRepo) UpdateTeacher(CourseID string, newTeacherID string) error {
	id, err := primitive.ObjectIDFromHex(CourseID)
	if err != nil {
		return err
	}

	teacherId, err := primitive.ObjectIDFromHex(newTeacherID)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"teacher", teacherId}}}}

	_, err = r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *CourseRepo) RemoveStudent(CourseID string, StudentID string, StudentRepo *StudentRepo) error {

	CourseObjectId, err := primitive.ObjectIDFromHex(CourseID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	Student, err := StudentRepo.FindStudentByID(StudentID)
	if err != nil {
		return fmt.Errorf("failed to find Student: %w", err)
	}
	if Student == nil {
		return fmt.Errorf("Student not found")
	}

	filter := bson.D{{"_id", CourseObjectId}}
	update := bson.D{{"$pull", bson.D{{"students", bson.D{{"$eq", Student}}}}}}

	_, err = r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete Student from Course: %w", err)
	}

	return nil
}

func (r *CourseRepo) GetCoursesByTeacherID(TeacherID string) ([]*types.Course, error) {

	ownerID, err := primitive.ObjectIDFromHex(TeacherID)
	if err != nil {
		return nil, fmt.Errorf("invalid TeacherID: %w", err)
	}

	filter := bson.M{"teacher": ownerID}
	var Courses []*types.Course

	cursor, err := r.MongoCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find Courses: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var Course types.Course
		if err := cursor.Decode(&Course); err != nil {
			return nil, fmt.Errorf("failed to decode Course: %w", err)
		}
		Courses = append(Courses, &Course)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return Courses, nil
}

func (r *CourseRepo) GetCoursesByStudentID(StudentID string) ([]*types.Course, error) {

	ownerID, err := primitive.ObjectIDFromHex(StudentID)
	if err != nil {
		return nil, fmt.Errorf("invalid StudentID: %w", err)
	}

	filter := bson.M{"students._id": ownerID}
	var Courses []*types.Course

	cursor, err := r.MongoCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find Courses: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var Course types.Course
		if err := cursor.Decode(&Course); err != nil {
			return nil, fmt.Errorf("failed to decode Course: %w", err)
		}
		Courses = append(Courses, &Course)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return Courses, nil
}
