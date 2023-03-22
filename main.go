package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type TeacherManager struct {
	mDBCache MDBCache
}

type Student struct {
	Code    string `json:"Code"`
	Name    string `json:"Name"`
	Program string `json:"Program"`
}

type conf struct {
	CassadraURI string `yaml:"cassandra-url"`
	CacheSize   int    `yaml:"cache-size"`
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
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

	var c conf
	c.getConf()

	mDBCache := LruInstantiate(c.CassadraURI, c.CacheSize)

	tMgr := TeacherManager{mDBCache: mDBCache}

	http.HandleFunc("/getstudents", tMgr.GetAllStudents)
	http.HandleFunc("/addstudent", tMgr.AddStudent)

	log.Println("Starting HTTP server...")

	http.ListenAndServe(":8080", nil)

}
