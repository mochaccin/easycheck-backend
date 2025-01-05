package handlers

import (
	"encoding/json"
	"money-minder/internal/database"
	"money-minder/internal/repositories"
	"money-minder/internal/types"
	"net/http"
)

var (
	service           = database.New()
	studentRepository = &repositories.StudentRepo{
		MongoCollection: service.GetCollection("students"),
	}
	courseRepository = &repositories.CourseRepo{
		MongoCollection: service.GetCollection("courses"),
	}
	attendanceRepository = &repositories.AttendanceRepo{
		MongoCollection: service.GetCollection("attendances"),
	}
)

func CreateStudent(w http.ResponseWriter, r *http.Request) error {
	usr := &types.Student{}
	derr := json.NewDecoder(r.Body).Decode(usr)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt create Student, verify that the values are formatted correctly",
		}
	}

	result, err := studentRepository.InsertStudent(usr)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, result)
}

func GetStudentByID(w http.ResponseWriter, r *http.Request) error {

	id := r.PathValue("id")

	Student, err := studentRepository.FindStudentByID(id)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, Student)
}

func AddStudentCourse(w http.ResponseWriter, r *http.Request) error {

	StudentId := r.PathValue("id")

	addCourseRequest := &StudentCourseRequest{}
	derr := json.NewDecoder(r.Body).Decode(addCourseRequest)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt add Course to Student, verify that the values are formatted correctly",
		}
	}

	err := studentRepository.AddCourse(StudentId, addCourseRequest.CourseId, courseRepository)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, "New Course added sucessfully.")
}

func RemoveStudentCourse(w http.ResponseWriter, r *http.Request) error {

	StudentId := r.PathValue("id")

	removeCourseRequest := &StudentCourseRequest{}
	derr := json.NewDecoder(r.Body).Decode(removeCourseRequest)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt remove Course to from, verify that the values are formatted correctly",
		}
	}

	err := studentRepository.RemoveCourse(StudentId, removeCourseRequest.CourseId, courseRepository)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, "Course deleted sucessfully.")
}

func AddStudentAttendance(w http.ResponseWriter, r *http.Request) error {

	StudentId := r.PathValue("id")

	addAttendanceRequest := &studentAttendanceRequest{}
	derr := json.NewDecoder(r.Body).Decode(addAttendanceRequest)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt add Attendance to Student, verify that the values are formatted correctly",
		}
	}

	err := studentRepository.AddAttendance(StudentId, addAttendanceRequest.AttendanceId, attendanceRepository)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, "New Attendance added sucessfully.")
}

func RemoveStudentAttendance(w http.ResponseWriter, r *http.Request) error {

	StudentId := r.PathValue("id")

	removeAttendanceRequest := &studentAttendanceRequest{}
	derr := json.NewDecoder(r.Body).Decode(removeAttendanceRequest)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt remove Attendance from Student, verify that the values are formatted correctly",
		}
	}

	err := studentRepository.RemoveAttendance(StudentId, removeAttendanceRequest.AttendanceId, attendanceRepository)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, "Attendance deleted sucessfully.")
}

func GetAllStudentsByCourseID(w http.ResponseWriter, r *http.Request) error {

	id := r.PathValue("id")

	Students, err := studentRepository.GetStudentsByCourseID(id)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, Students)
}

type StudentCourseRequest struct {
	CourseId string `json:"CourseId" bson:"course_id"`
}

type studentAttendanceRequest struct {
	AttendanceId string `json:"AttendanceId" bson:"attendance_id"`
}
