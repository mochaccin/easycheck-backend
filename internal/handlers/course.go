package handlers

import (
	"encoding/json"
	"money-minder/internal/types"
	"net/http"
)

func CreateCourse(w http.ResponseWriter, r *http.Request) error {
	Course := &types.Course{}
	derr := json.NewDecoder(r.Body).Decode(Course)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Could not create Course, verify that the values are formatted correctly",
		}
	}

	result, err := courseRepository.InsertCourse(Course)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, result)
}

func DeleteCourse(w http.ResponseWriter, r *http.Request) error {

	CourseId := r.PathValue("id")

	result, err := courseRepository.DeleteCourse(CourseId)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, result)
}

func GetCourseByID(w http.ResponseWriter, r *http.Request) error {

	id := r.PathValue("id")

	Course, err := courseRepository.FindCourseByID(id)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, Course)
}

func GetAllCoursesByStudentID(w http.ResponseWriter, r *http.Request) error {

	id := r.PathValue("id")

	Courses, err := courseRepository.GetCoursesByStudentID(id)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, Courses)
}

func GetAllCoursesByTeacherID(w http.ResponseWriter, r *http.Request) error {

	id := r.PathValue("id")

	Courses, err := courseRepository.GetCoursesByTeacherID(id)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, Courses)
}

func AddCourseStudent(w http.ResponseWriter, r *http.Request) error {

	CourseId := r.PathValue("id")

	addStudentRequest := &CourseStudentRequest{}
	derr := json.NewDecoder(r.Body).Decode(addStudentRequest)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt add Student to Course, verify that the values are formatted correctly",
		}
	}

	err := courseRepository.AddStudent(CourseId, addStudentRequest.StudentId, studentRepository)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, "New Student added sucessfully.")
}

func RemoveCourseStudent(w http.ResponseWriter, r *http.Request) error {

	CourseId := r.PathValue("id")

	removeStudentRequest := &CourseStudentRequest{}
	derr := json.NewDecoder(r.Body).Decode(removeStudentRequest)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt remove Student from Course, verify that the values are formatted correctly",
		}
	}
	err := courseRepository.RemoveStudent(CourseId, removeStudentRequest.StudentId, studentRepository)

	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, "Student deleted sucessfully.")
}

func UpdateCourseTeacher(w http.ResponseWriter, r *http.Request) error {

	CourseId := r.PathValue("id")

	addTeacherRequest := &CourseTeacherRequest{}
	derr := json.NewDecoder(r.Body).Decode(addTeacherRequest)

	if derr != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Couldnt update Teacher to Course, verify that the values are formatted correctly",
		}
	}

	err := courseRepository.UpdateTeacher(CourseId, addTeacherRequest.TeacherId)
	if err != nil {
		return APIError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return WriteJSON(w, http.StatusOK, "New Course Teacher updated sucessfully.")
}

type CourseStudentRequest struct {
	StudentId string `json:"StudentId" bson:"student_id"`
}

type CourseTeacherRequest struct {
	TeacherId string `json:"teacherId" bson:"teacher_id"`
}
