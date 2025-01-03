package repositories

import (
	"context"
	"money-minder/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AttendanceRepo struct {
	MongoCollection *mongo.Collection
}

func (r *AttendanceRepo) InsertAttendance(attendance *types.Attendance) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), attendance)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *AttendanceRepo) DeleteAttendance(attendanceID string) error {
	id, err := primitive.ObjectIDFromHex(attendanceID)
	if err != nil {
		return err
	}
	_, err = r.MongoCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *AttendanceRepo) FindAttendanceByID(attendanceID string) (*types.Attendance, error) {
	id, err := primitive.ObjectIDFromHex(attendanceID)
	if err != nil {
		return nil, err
	}
	var attendance types.Attendance
	err = r.MongoCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&attendance)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &attendance, nil
}

func (r *AttendanceRepo) UpdateAttendance(attendanceID string, updatedAttendance *types.Attendance) error {
	id, err := primitive.ObjectIDFromHex(attendanceID)
	if err != nil {
		return err
	}
	_, err = r.MongoCollection.ReplaceOne(context.Background(), bson.M{"_id": id}, updatedAttendance)
	return err
}

func (r *AttendanceRepo) GetAllAttendanceByCourseID(courseID string) ([]types.Attendance, error) {
	id, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.MongoCollection.Find(context.Background(), bson.M{"course": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var attendances []types.Attendance
	if err = cursor.All(context.Background(), &attendances); err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *AttendanceRepo) GetAllAttendanceByStudentID(studentID string) ([]types.Attendance, error) {
	id, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.MongoCollection.Find(context.Background(), bson.M{"students": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var attendances []types.Attendance
	if err = cursor.All(context.Background(), &attendances); err != nil {
		return nil, err
	}
	return attendances, nil
}
