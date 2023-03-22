package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type TeacherManager struct {
	mDBCache MDBCache
}

type Student struct {
	Code    string `json:"Code"`
	Name    string `json:"Name"`
	Program string `json:"Program"`
}

func (t *TeacherManager) GetAllStudents(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	teacherId, present := query["teacherId"]
	if !present || len(teacherId) == 0 {
		log.Fatal("Please provide teacherId")
	}

	teacherMDB := t.mDBCache.GetMDB(teacherId[0])

	CreateTableIfNotExists(teacherMDB)

	students := GetAllStudents(teacherMDB)

	json.NewEncoder(w).Encode(students)

}
func (t *TeacherManager) AddStudent(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	teacherId, present := query["teacherId"]
	if !present || len(teacherId) == 0 {
		log.Fatal("Please provide teacherId")
	}

	teacherMDB := t.mDBCache.GetMDB(teacherId[0])

	CreateTableIfNotExists(teacherMDB)
	var s Student

	json.NewDecoder(r.Body).Decode(&s)

	AddStudent(teacherMDB, s)

	json.NewEncoder(w).Encode(s)

}

func main() {

	mDBCache := LruInstantiate()

	c := TeacherManager{mDBCache: mDBCache}

	http.HandleFunc("/getstudents", c.GetAllStudents)
	http.HandleFunc("/addstudent", c.AddStudent)

	log.Println("Starting HTTP server...")

	http.ListenAndServe(":8080", nil)

}
