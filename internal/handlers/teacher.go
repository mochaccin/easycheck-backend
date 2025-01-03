package handlers

import (
	"encoding/json"
	"money-minder/internal/auth"
	"money-minder/internal/repositories"
	"money-minder/internal/types"
	"net/http"
)

var (
	teacherRepo = &repositories.TeacherRepo{
		MongoCollection: service.GetCollection("teachers"),
	}
)

func CreateTeacher(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can create other teachers"}
	}

	var teacher types.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}

	result, err := teacherRepo.InsertTeacher(&teacher)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}

	return WriteJSON(w, http.StatusCreated, result)
}

func DeleteTeacher(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can delete other teachers"}
	}

	teacherID := r.PathValue("teacherID")
	if err := teacherRepo.DeleteTeacher(teacherID); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Teacher deleted successfully")
}

func GetTeacherByID(w http.ResponseWriter, r *http.Request) error {
	_, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	teacherID := r.PathValue("teacherID")
	teacher, err := teacherRepo.FindTeacherByID(teacherID)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	if teacher == nil {
		return APIError{Status: http.StatusNotFound, Msg: "Teacher not found"}
	}
	return WriteJSON(w, http.StatusOK, teacher)
}

func GetAllTeachers(w http.ResponseWriter, r *http.Request) error {
	_, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	teachers, err := teacherRepo.FindAllTeachers()
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, teachers)
}

func AddCourseToTeacher(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can add courses to teachers"}
	}

	teacherID := r.PathValue("teacherID")
	var course types.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}
	if err := teacherRepo.AddCourse(teacherID, course); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Course added to teacher successfully")
}

func RemoveCourseFromTeacher(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can remove courses from teachers"}
	}

	teacherID := r.PathValue("teacherID")
	var req struct {
		CourseID string `json:"courseId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}
	if err := teacherRepo.RemoveCourse(teacherID, req.CourseID); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Course removed from teacher successfully")
}

func GetAllCoursesByTeacherID(w http.ResponseWriter, r *http.Request) error {
	_, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	teacherID := r.PathValue("teacherID")
	courses, err := teacherRepo.GetAllCoursesByTeacherID(teacherID)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, courses)
}
