package handlers

import (
	"encoding/json"
	"money-minder/internal/types"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

type Claims struct {
	ID    interface{} `json:"id"`
	Email string      `json:"email"`
	Role  string      `json:"role"`
	jwt.RegisteredClaims
}

func Register(w http.ResponseWriter, r *http.Request) error {
	var registerRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"` // "student" or "teacher"
	}

	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: "Error hashing password"}
	}

	var result interface{}

	switch registerRequest.Role {
	case "student":
		student := &types.Student{
			Name:     registerRequest.Name,
			Email:    registerRequest.Email,
			Password: string(hashedPassword),
		}
		result, err = studentRepo.InsertStudent(student)
	case "teacher":
		teacher := &types.Teacher{
			Name:     registerRequest.Name,
			Email:    registerRequest.Email,
			Password: string(hashedPassword),
		}
		result, err = teacherRepo.InsertTeacher(teacher)
	default:
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid role"}
	}

	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}

	return WriteJSON(w, http.StatusCreated, result)
}

func Login(w http.ResponseWriter, r *http.Request) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"` // "student" or "teacher"
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid request body"}
	}

	var (
		id       interface{}
		password string
	)

	switch loginRequest.Role {
	case "student":
		student, err := studentRepo.FindStudentByEmail(loginRequest.Email)
		if err != nil {
			return APIError{Status: http.StatusUnauthorized, Msg: "Invalid credentials"}
		}
		id = student.ID
		password = student.Password
	case "teacher":
		teacher, err := teacherRepo.FindTeacherByEmail(loginRequest.Email)
		if err != nil {
			return APIError{Status: http.StatusUnauthorized, Msg: "Invalid credentials"}
		}
		id = teacher.ID
		password = teacher.Password
	default:
		return APIError{Status: http.StatusBadRequest, Msg: "Invalid role"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(loginRequest.Password)); err != nil {
		return APIError{Status: http.StatusUnauthorized, Msg: "Invalid credentials"}
	}

	// Create claims with expiry
	claims := &Claims{
		ID:    id,
		Email: loginRequest.Email,
		Role:  loginRequest.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return APIError{Status: http.StatusInternalServerError, Msg: "Error generating token"}
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}
