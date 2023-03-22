package main

import (
	"database/sql"
	"log"
)

func CreateTableIfNotExists(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE IF NOT EXISTS student (
		"idStudent" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"code" TEXT,
		"name" TEXT,
		"program" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create student table if not exists...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	//log.Println("student table created")
}

func GetAllStudents(db *sql.DB) []Student {
	row, err := db.Query("SELECT * FROM student ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	students := []Student{}

	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var code string
		var name string
		var program string
		row.Scan(&id, &code, &name, &program)
		students = append(students, Student{code, name, program})
		//log.Println("Student: ", code, " ", name, " ", program)
	}

	return students
}

func AddStudent(db *sql.DB, s Student) {
	//log.Println("Inserting student record ...")
	insertStudentSQL := `INSERT INTO student(code, name, program) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(s.Code, s.Name, s.Program)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
