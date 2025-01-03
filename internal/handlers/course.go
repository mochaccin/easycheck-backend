package handlers

import (
	"encoding/json"
	"money-minder/internal/auth"
	"money-minder/internal/repositories"
	"money-minder/internal/types"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	courseRepo = &repositories.CourseRepo{
		MongoCollection: service.GetCollection("courses"),
	}
)

func CreateCourse(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can create courses"}
	}

	var course types.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}

	// Set the teacher ID from the JWT token
	teacherID, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid teacher ID"}
	}
	course.Teacher = teacherID

	result, err := courseRepo.InsertCourse(&course)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}

	return WriteJSON(w, http.StatusCreated, result)
}

func DeleteCourse(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can delete courses"}
	}

	courseID := r.PathValue("courseID")
	if err := courseRepo.DeleteCourse(courseID); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Course deleted successfully")
}

func GetCourseByID(w http.ResponseWriter, r *http.Request) error {
	_, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	courseID := r.PathValue("courseID")
	course, err := courseRepo.FindCourseByID(courseID)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	if course == nil {
		return APIError{Status: http.StatusNotFound, Msg: "Course not found"}
	}
	return WriteJSON(w, http.StatusOK, course)
}

func ChangeTeacher(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can change course teachers"}
	}

	courseID := r.PathValue("courseID")
	var req struct {
		TeacherID string `json:"teacherId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}
	if err := courseRepo.ChangeTeacher(courseID, req.TeacherID); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Teacher changed successfully")
}

func AddStudentToCourse(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can add students to courses"}
	}

	courseID := r.PathValue("courseID")
	var student types.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}
	if err := courseRepo.AddStudent(courseID, student); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Student added to course successfully")
}

func RemoveStudentFromCourse(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can remove students from courses"}
	}

	courseID := r.PathValue("courseID")
	var req struct {
		StudentID string `json:"studentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}
	if err := courseRepo.RemoveStudent(courseID, req.StudentID); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Student removed from course successfully")
}

func GetAllStudentsByCourseID(w http.ResponseWriter, r *http.Request) error {
	_, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	courseID := r.PathValue("courseID")
	students, err := courseRepo.GetAllStudentsByCourseID(courseID)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, students)
}

func GetAllCourses(w http.ResponseWriter, r *http.Request) error {
	_, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	courses, err := courseRepo.FindAllCourses()
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, courses)
}
