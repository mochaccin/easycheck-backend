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
	MongoCollection *mongo.Collection
}

func (r *StudentRepo) InsertStudent(usr *types.Student) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne((context.Background()), usr)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *StudentRepo) DeleteStudent(usr *types.Student) (interface{}, error) {
	result, err := r.MongoCollection.DeleteOne((context.Background()), usr)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *StudentRepo) FindStudentByID(usrID string) (*types.Student, error) {
	id, err := primitive.ObjectIDFromHex(usrID)

	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}

	var usr types.Student

	err = r.MongoCollection.FindOne(context.Background(), filter).Decode(&usr)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &usr, nil
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

func (r *StudentRepo) FindAllStudents() ([]types.Student, error) {
	results, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var usrs []types.Student

	err = results.All(context.Background(), &usrs)
	if err != nil {
		return nil, fmt.Errorf("Find all uses results decode error %s", err.Error())
	}

	return usrs, nil
}

func (r *StudentRepo) AddCourse(StudentID string, CourseID string, CourseRepo *CourseRepo) error {

	StudentObjectID, err := primitive.ObjectIDFromHex(StudentID)
	if err != nil {
		return fmt.Errorf("invalid Student ID: %w", err)
	}

	Course, err := CourseRepo.FindCourseByID(CourseID)
	if err != nil {
		return fmt.Errorf("failed to find Course: %w", err)
	}
	if Course == nil {
		return fmt.Errorf("Course not found")
	}

	filter := bson.D{{"_id", StudentObjectID}}
	update := bson.D{{"$push", bson.D{{"courses", Course}}}}

	_, err = r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to add Course to Student: %w", err)
	}

	return nil
}

func (r *StudentRepo) RemoveCourse(StudentID string, CourseID string, CourseRepo *CourseRepo) error {

	StudentObjectID, err := primitive.ObjectIDFromHex(StudentID)
	if err != nil {
		return fmt.Errorf("invalid Student ID: %w", err)
	}

	Course, err := CourseRepo.FindCourseByID(CourseID)
	if err != nil {
		return fmt.Errorf("failed to find Course: %w", err)
	}
	if Course == nil {
		return fmt.Errorf("Course not found")
	}

	filter := bson.D{{"_id", StudentObjectID}}
	update := bson.D{{"$pull", bson.D{{"courses", bson.D{{"$eq", Course}}}}}}

	_, err = r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete Course from Student: %w", err)
	}

	return nil
}

func (r *StudentRepo) AddAttendance(StudentID string, AttendanceID string, AttendanceRepo *AttendanceRepo) error {

	StudentObjID, err := primitive.ObjectIDFromHex(StudentID)
	if err != nil {
		return fmt.Errorf("invalid Student ID: %w", err)
	}

	Attendance, err := AttendanceRepo.FindAttendanceByID(AttendanceID)
	if err != nil {
		return fmt.Errorf("failed to find Attendance: %w", err)
	}
	if Attendance == nil {
		return fmt.Errorf("Attendance not found")
	}

	filter := bson.D{{"_id", StudentObjID}}
	update := bson.D{{"$push", bson.D{{"Attendances", Attendance}}}}

	_, err = r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to add Attendance to Student: %w", err)
	}

	return nil
}

func (r *StudentRepo) RemoveAttendance(StudentID string, AttendanceID string, AttendanceRepo *AttendanceRepo) error {

	StudentObjectID, err := primitive.ObjectIDFromHex(StudentID)
	if err != nil {
		return fmt.Errorf("invalid Student ID: %w", err)
	}

	Attendance, err := AttendanceRepo.FindAttendanceByID(AttendanceID)
	if err != nil {
		return fmt.Errorf("failed to find Attendance: %w", err)
	}
	if Attendance == nil {
		return fmt.Errorf("Attendance not found")
	}

	filter := bson.D{{"_id", StudentObjectID}}
	update := bson.D{{"$pull", bson.D{{"Attendances", bson.D{{"$eq", Attendance}}}}}}

	_, err = r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete Attendance from Student: %w", err)
	}

	return nil
}

func (r *StudentRepo) GetStudentsByCourseID(CourseID string) ([]*types.Student, error) {

	courseID, err := primitive.ObjectIDFromHex(CourseID)
	if err != nil {
		return nil, fmt.Errorf("invalid courseID: %w", err)
	}

	filter := bson.M{"course_id": courseID}
	var Students []*types.Student

	cursor, err := r.MongoCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find Students: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var Student types.Student
		if err := cursor.Decode(&Student); err != nil {
			return nil, fmt.Errorf("failed to decode Attendance: %w", err)
		}
		Students = append(Students, &Student)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return Students, nil
}
