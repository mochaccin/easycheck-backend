package server

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"money-minder/internal/auth"
	"money-minder/internal/handlers"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.healthHandler)

	// Student routes
	mux.HandleFunc("POST /students", makeHandler(handlers.CreateStudent))
	mux.HandleFunc("GET /students/{studentID}", makeHandler(handlers.GetStudentByID))
	mux.HandleFunc("PATCH /students/{studentID}/courses", makeHandler(handlers.AddStudentCourse))
	mux.HandleFunc("DELETE /students/{studentID}/courses", makeHandler(handlers.RemoveStudentCourse))
	mux.HandleFunc("GET /students/{studentID}/courses", jwtMiddleware(makeHandler(handlers.GetAllCoursesByStudentID)))
	mux.HandleFunc("GET /students/{studentID}/attendances", jwtMiddleware(makeHandler(handlers.GetAllAttendancesByStudentID)))

	// Teacher routes
	mux.HandleFunc("POST /teachers", makeHandler(handlers.CreateTeacher))
	mux.HandleFunc("GET /teachers/{teacherID}", makeHandler(handlers.GetTeacherByID))
	mux.HandleFunc("PATCH /teachers/{teacherID}/courses", makeHandler(handlers.AddTeacherCourse))
	mux.HandleFunc("DELETE /teachers/{teacherID}/courses", makeHandler(handlers.RemoveTeacherCourse))
	mux.HandleFunc("GET /teachers/{teacherID}/courses", jwtMiddleware(makeHandler(handlers.GetAllCoursesByTeacherID)))

	// Course routes
	mux.HandleFunc("POST /courses", makeHandler(handlers.CreateCourse))
	mux.HandleFunc("GET /courses/{courseID}", makeHandler(handlers.GetCourseByID))
	mux.HandleFunc("DELETE /courses/{courseID}", makeHandler(handlers.DeleteCourse))
	mux.HandleFunc("PATCH /courses/{courseID}/teacher", makeHandler(handlers.UpdateCourseTeacher))
	mux.HandleFunc("PATCH /courses/{courseID}/students", makeHandler(handlers.AddCourseStudent))
	mux.HandleFunc("DELETE /courses/{courseID}/students", makeHandler(handlers.RemoveCourseStudent))
	mux.HandleFunc("GET /courses/{courseID}/students", jwtMiddleware(makeHandler(handlers.GetAllStudentsByCourseID)))

	// Attendance routes
	mux.HandleFunc("POST /attendance", makeHandler(handlers.CreateAttendance))
	mux.HandleFunc("PATCH /attendance/{attendanceID}", makeHandler(handlers.UpdateAttendance))
	mux.HandleFunc("DELETE /attendance/{attendanceID}", makeHandler(handlers.DeleteAttendance))
	mux.HandleFunc("GET /attendance/byCourse/{courseID}", jwtMiddleware(makeHandler(handlers.GetAllAttendancesByCourseID)))
	mux.HandleFunc("GET /attendance/byStudent/{studentID}", jwtMiddleware(makeHandler(handlers.GetAllAttendancesByStudentID)))

	// Auth routes
	mux.HandleFunc("POST /auth/register", makeHandler(handlers.Register))
	mux.HandleFunc("POST /auth/login", makeHandler(handlers.Login))

	return corsMiddleware(mux)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHandler(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if e, ok := err.(handlers.APIError); ok {
				slog.Error("API Error", "err", e, "status", e.Status)
				handlers.WriteJSON(w, e.Status, e)
			}
		}
	}
}

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Check if it starts with "Bearer "
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, prefix)

		// Parse and validate the token
		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			slog.Error("JWT Parse error", "error", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Token is not valid", http.StatusUnauthorized)
			return
		}

		// Add claims to the request context
		ctx := auth.SetClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
