package handlers

import (
	"encoding/json"
	"money-minder/internal/types"
	"net/http"
)

func CreateAttendance(w http.ResponseWriter, r *http.Request) error {
	Attendance := &types.Attendance{}
	derr := json.NewDecoder(r.Body).Decode(Attendance)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Could not create Attendance, verify that the values are formatted correctly",
		}
	}

	result, err := attendanceRepository.InsertAttendance(Attendance)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, result)
}

func DeleteAttendance(w http.ResponseWriter, r *http.Request) error {

	AttendanceId := r.PathValue("id")

	result, err := attendanceRepository.DeleteAttendance(AttendanceId)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, result)
}

func UpdateAttendance(w http.ResponseWriter, r *http.Request) error {

	userId := r.PathValue("id")

	updateUsr := &UpdateAttendanceRequest{}
	derr := json.NewDecoder(r.Body).Decode(updateUsr)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt update attendance isPresent value, verify that the values are formatted correctly",
		}
	}

	err := attendanceRepository.UpdatePresent(userId, updateUsr.IsPresent)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, "Attendance's isPresent value updated succesfully")
}

func GetAttendanceByID(w http.ResponseWriter, r *http.Request) error {

	AttendanceId := r.PathValue("id")

	Attendance, err := attendanceRepository.FindAttendanceByID(AttendanceId)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, Attendance)
}

func GetAllAttendancesByCourseID(w http.ResponseWriter, r *http.Request) error {

	CourseID := r.PathValue("id")

	cards, err := attendanceRepository.GetAttendancesByCourseID(CourseID)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, cards)
}

func GetAllAttendancesByStudentID(w http.ResponseWriter, r *http.Request) error {

	StudentID := r.PathValue("id")

	attendances, err := attendanceRepository.GetAttendancesByStudentID(StudentID)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, attendances)
}

type UpdateAttendanceRequest struct {
	IsPresent bool `json:"isPresent" bson:"present"`
}
