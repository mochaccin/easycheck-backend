# Easycheck

Si alguna wea falla taggenme en el ds

## Prerequisitos

- Go
- Configurar el .env

## Como instalar go

- Descargar la version que corresponda a su sistema operativo desde aca [Go Downloads](https://go.dev/dl/)
- Lo instalan como cualquier otro programa

Listo ya tienen go instalado.

Ahora clonan el repo y tienen 2 opciones, buildearlo manual de la siguiente forma

Descargar dependencias

```go
go mod tidy
```
Correr el proyecto

```go
go run cmd/api/main.go
```

La otra opcion es utilizar live reloading, cosa que si quieren cambiar algo no tienen que reiniciar como tal.

Para eso hay que usar go/air

Instalan chocolatey que nos permite instalar herramientas de linux en windows, en este caso el comando make para usar el makefile que usa go air.

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

Despues instalan go air para el hot reloading

```powershell
go install github.com/air-verse/air@latest
```

Finalmente para ejecutar el proyecto abren una terminal en el vscode y escriben

```powershell
air
```

## Endpoints

| Resource | HTTP Method | Endpoint | Request Body | Response
|-----|-----|-----|-----|-----
| **Student**
| Create Student | POST | /students | Student object | Created student object
| Get Student by ID | GET | /students/{studentID} | - | Student object
| Delete Student | DELETE | /students/{studentID} | - | Success message
| Add Course to Student | PATCH | /students/{studentID}/courses | Course object | Success message
| Remove Course from Student | DELETE | /students/{studentID}/courses | { "courseId": "string" } | Success message
| Get All Courses by Student ID | GET | /students/{studentID}/courses | - | Array of Course objects
| Get All Attendances by Student ID | GET | /students/{studentID}/attendances | - | Array of Attendance objects
| Get All Students | GET | /students | - | Array of Student objects
| **Teacher**
| Create Teacher | POST | /teachers | Teacher object | Created teacher object
| Get Teacher by ID | GET | /teachers/{teacherID} | - | Teacher object
| Delete Teacher | DELETE | /teachers/{teacherID} | - | Success message
| Add Course to Teacher | PATCH | /teachers/{teacherID}/courses | Course object | Success message
| Remove Course from Teacher | DELETE | /teachers/{teacherID}/courses | { "courseId": "string" } | Success message
| Get All Courses by Teacher ID | GET | /teachers/{teacherID}/courses | - | Array of Course objects
| Get All Teachers | GET | /teachers | - | Array of Teacher objects
| **Course**
| Create Course | POST | /courses | Course object | Created course object
| Get Course by ID | GET | /courses/{courseID} | - | Course object
| Delete Course | DELETE | /courses/{courseID} | - | Success message
| Change Teacher | PATCH | /courses/{courseID}/teacher | { "teacherId": "string" } | Success message
| Add Student to Course | PATCH | /courses/{courseID}/students | Student object | Success message
| Remove Student from Course | DELETE | /courses/{courseID}/students | { "studentId": "string" } | Success message
| Get All Students by Course ID | GET | /courses/{courseID}/students | - | Array of Student objects
| **Attendance**
| Create Attendance | POST | /attendance | Attendance object | Created attendance object
| Update Attendance | PATCH | /attendance/{attendanceID} | Updated Attendance object | Success message
| Delete Attendance | DELETE | /attendance/{attendanceID} | - | Success message
| Get All Attendance by Course ID | GET | /attendance/byCourse/{courseID} | - | Array of Attendance objects
| Get All Attendance by Student ID | GET | /attendance/byStudent/{studentID} | - | Array of Attendance objects
| **Auth**
| Register | POST | /auth/register | Registration details | JWT token
| Login | POST | /auth/login | Login credentials | JWT token
