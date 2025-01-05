package handlers

import (
	"encoding/json"
	"money-minder/internal/repositories"
	"money-minder/internal/types"
	"net/http"
)

var (
	teacherRepository = &repositories.TeacherRepo{
		MongoCollection: service.GetCollection("teachers"),
	}
)

func CreateTeacher(w http.ResponseWriter, r *http.Request) error {
	usr := &types.Teacher{}
	derr := json.NewDecoder(r.Body).Decode(usr)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt create Teacher, verify that the values are formatted correctly",
		}
	}

	result, err := teacherRepository.InsertTeacher(usr)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, result)
}

func GetTeacherByID(w http.ResponseWriter, r *http.Request) error {

	id := r.PathValue("id")

	Teacher, err := teacherRepository.FindTeacherByID(id)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, Teacher)
}

func AddTeacherCourse(w http.ResponseWriter, r *http.Request) error {

	TeacherId := r.PathValue("id")

	addCourseRequest := &TeacherCourseRequest{}
	derr := json.NewDecoder(r.Body).Decode(addCourseRequest)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt add Course to Teacher, verify that the values are formatted correctly",
		}
	}

	err := teacherRepository.AddCourse(TeacherId, addCourseRequest.CourseId, courseRepository)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, "New Course added sucessfully.")
}

func RemoveTeacherCourse(w http.ResponseWriter, r *http.Request) error {

	TeacherId := r.PathValue("id")

	removeCourseRequest := &TeacherCourseRequest{}
	derr := json.NewDecoder(r.Body).Decode(removeCourseRequest)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt remove Course to from, verify that the values are formatted correctly",
		}
	}

	err := teacherRepository.RemoveCourse(TeacherId, removeCourseRequest.CourseId, courseRepository)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, "Course deleted sucessfully.")
}

type TeacherCourseRequest struct {
	CourseId string `json:"CourseId" bson:"course_id"`
}

type TeacherAttendanceRequest struct {
	AttendanceId string `json:"AttendanceId" bson:"attendance_id"`
}
