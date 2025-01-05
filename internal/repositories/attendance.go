package repositories

import (
	"context"
	"fmt"
	"money-minder/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AttendanceRepo struct {
	MongoCollection *mongo.Collection
}

func (r *AttendanceRepo) InsertAttendance(Attendance *types.Attendance) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne((context.Background()), Attendance)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AttendanceRepo) DeleteAttendance(AttendanceID string) (*mongo.DeleteResult, error) {
	id, err := primitive.ObjectIDFromHex(AttendanceID)
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

func (r *AttendanceRepo) UpdatePresent(usrID string, isPresent bool) error {
	id, err := primitive.ObjectIDFromHex(usrID)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"present", isPresent}}}}

	_, err = r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *AttendanceRepo) FindAttendanceByID(AttendanceID string) (*types.Attendance, error) {
	id, err := primitive.ObjectIDFromHex(AttendanceID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}

	var Attendance types.Attendance

	err = r.MongoCollection.FindOne(context.Background(), filter).Decode(&Attendance)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &Attendance, nil
}

func (r *AttendanceRepo) GetAttendancesByCourseID(id string) ([]*types.Attendance, error) {

	CourseID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid CourseID: %w", err)
	}

	filter := bson.M{"course_id": CourseID}
	var Attendances []*types.Attendance

	cursor, err := r.MongoCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find Attendances: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var Attendance types.Attendance
		if err := cursor.Decode(&Attendance); err != nil {
			return nil, fmt.Errorf("failed to decode Attendance: %w", err)
		}
		Attendances = append(Attendances, &Attendance)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return Attendances, nil
}

func (r *AttendanceRepo) GetAttendancesByStudentID(StudentID string) ([]*types.Attendance, error) {

	ownerID, err := primitive.ObjectIDFromHex(StudentID)
	if err != nil {
		return nil, fmt.Errorf("invalid StudentID: %w", err)
	}

	filter := bson.M{"student_id": ownerID}
	var Attendances []*types.Attendance

	cursor, err := r.MongoCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find Attendances: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var Attendance types.Attendance
		if err := cursor.Decode(&Attendance); err != nil {
			return nil, fmt.Errorf("failed to decode Attendance: %w", err)
		}
		Attendances = append(Attendances, &Attendance)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return Attendances, nil
}
