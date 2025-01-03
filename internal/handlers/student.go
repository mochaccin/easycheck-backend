package handlers

import (
	"encoding/json"
	"money-minder/internal/auth"
	"money-minder/internal/database"
	"money-minder/internal/repositories"
	"money-minder/internal/types"
	"net/http"
)

var (
	service     = database.New()
	studentRepo = &repositories.StudentRepo{
		MongoCollection: service.GetCollection("students"),
	}
)

func CreateStudent(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can create students"}
	}

	var student types.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}

	result, err := studentRepo.InsertStudent(&student)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}

	return WriteJSON(w, http.StatusCreated, result)
}

func DeleteStudent(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can delete students"}
	}

	studentID := r.PathValue("studentID")
	if err := studentRepo.DeleteStudent(studentID); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Student deleted successfully")
}

func GetStudentByID(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	studentID := r.PathValue("studentID")
	if claims.Role == "student" && claims.ID != studentID {
		return APIError{Status: http.StatusForbidden, Msg: "Students can only view their own information"}
	}

	student, err := studentRepo.FindStudentByID(studentID)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	if student == nil {
		return APIError{Status: http.StatusNotFound, Msg: "Student not found"}
	}
	return WriteJSON(w, http.StatusOK, student)
}

func GetAllStudents(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can view all students"}
	}

	students, err := studentRepo.FindAllStudents()
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, students)
}

func AddCourseToStudent(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can add courses to students"}
	}

	studentID := r.PathValue("studentID")
	var course types.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}
	if err := studentRepo.AddCourse(studentID, course); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Course added to student successfully")
}

func RemoveCourseFromStudent(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can remove courses from students"}
	}

	studentID := r.PathValue("studentID")
	var req struct {
		CourseID string `json:"courseId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}
	if err := studentRepo.RemoveCourse(studentID, req.CourseID); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Course removed from student successfully")
}

func GetAllCoursesByStudentID(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	studentID := r.PathValue("studentID")
	if claims.Role == "student" && claims.ID != studentID {
		return APIError{Status: http.StatusForbidden, Msg: "Students can only view their own courses"}
	}

	courses, err := studentRepo.GetAllCoursesByStudentID(studentID)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, courses)
}

func GetAllAttendancesByStudentID(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	studentID := r.PathValue("studentID")
	if claims.Role == "student" && claims.ID != studentID {
		return APIError{Status: http.StatusForbidden, Msg: "Students can only view their own attendance"}
	}

	attendances, err := studentRepo.GetAllAttendancesByStudentID(studentID)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, attendances)
}
