package handlers

import (
	"encoding/json"
	"money-minder/internal/auth"
	"money-minder/internal/repositories"
	"money-minder/internal/types"
	"net/http"
)

var (
	attendanceRepo = &repositories.AttendanceRepo{
		MongoCollection: service.GetCollection("attendances"),
	}
)

func CreateAttendance(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can create attendance"}
	}

	var attendance types.Attendance
	if err := json.NewDecoder(r.Body).Decode(&attendance); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}

	result, err := attendanceRepo.InsertAttendance(&attendance)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}

	return WriteJSON(w, http.StatusCreated, result)
}

func DeleteAttendance(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can delete attendance"}
	}

	attendanceID := r.PathValue("attendanceID")
	if err := attendanceRepo.DeleteAttendance(attendanceID); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Attendance deleted successfully")
}

func UpdateAttendance(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can update attendance"}
	}

	attendanceID := r.PathValue("attendanceID")
	var updatedAttendance types.Attendance
	if err := json.NewDecoder(r.Body).Decode(&updatedAttendance); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}
	if err := attendanceRepo.UpdateAttendance(attendanceID, &updatedAttendance); err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, "Attendance updated successfully")
}

func GetAllAttendanceByCourseID(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if claims.Role != "teacher" {
		return APIError{Status: http.StatusForbidden, Msg: "Only teachers can see the course attendance"}
	}

	courseID := r.PathValue("courseID")
	attendances, err := attendanceRepo.GetAllAttendanceByCourseID(courseID)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, attendances)
}

func GetAllAttendanceByStudentID(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	studentID := r.PathValue("studentID")
	if claims.Role == "student" && claims.ID != studentID {
		return APIError{Status: http.StatusForbidden, Msg: "Students can only view their own attendance"}
	}

	attendances, err := attendanceRepo.GetAllAttendanceByStudentID(studentID)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return WriteJSON(w, http.StatusOK, attendances)
}
