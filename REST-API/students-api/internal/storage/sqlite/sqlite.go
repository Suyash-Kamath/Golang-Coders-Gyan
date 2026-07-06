package sqlite

// implementing the interface inside this

// inbuilt package
// database se kaam karnekeliye jo bhi interface hai wo provide karta hai
import (
	"database/sql"
	"fmt"

	"github.com/Suyash-Kamath/Golang-Coders-Gyan/REST-API/students-api/internal/config"
	"github.com/Suyash-Kamath/Golang-Coders-Gyan/REST-API/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
	// Dekho, we are not using driver ka variable , bass insitialisation keliye chahiye hota hai , we are using it indirectly behind the scene , so underscore lagaayaa
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	// Open database , means connection
	// database ka driver pass and storage path
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)
	// hey _,err:= me why did you remove the colon

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students(name,email,age) VALUES (?,?,?)")
	// we dont pass the data directly , first we prepare the query and give the placeholder and then we bind the data , which prevenets the SQL injection attack

	if err != nil {
		return 0, err
	}
	// stmt ko bhi close karna hota hai
	defer stmt.Close()

	// as we have prepared the query so execute bhi karna hai

	result, err := stmt.Exec(name, email, age) // binds the stmt and replaces the placeholder

	if err != nil {
		return 0, err
	}

	// result ke andhar , seeing the documentation , we get LastInsertID() , RowsAffected() function

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
		// int64 ka empty value is 0
	}
	return lastId, nil

}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {

	// stmt,err:=s.Db.Prepare(`SELECT * FROM students WHERE id = ? LIMIT 1`)
	stmt, err := s.Db.Prepare(`SELECT id,name,email,age FROM students WHERE id = ? LIMIT 1`)
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	// data jo database se , aar aha hai , struct ke andhar serialize karke daaldo

	var student types.Student
	// Scan method database ka data , struct me dalega , kya kya scan karna hai wo batana hai , and order wise hona chahiye .
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		// if no user found
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("No student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("Query error %w", err)
	}

	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	// Always use paginations for such queries
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students ")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	// yahape koi placeholder nahi hai toh , pass karne ki jarurat nahi
	rows, err := stmt.Query() // returns rows and error
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student

	//rows yeh list hai , so loop
	// Rows is a result of query . Its cursor starts before the first row
	//  of the result set. Use rows.Next() to advance from row to row

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) UpdateStudentById(id int64, name string, email string, age int) error {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No student found with id %d", id)
	}

	return nil

}

func (s *Sqlite) DeleteStudentById(id int64) error {
	stmt, err := s.Db.Prepare("Delete from students where id = ?")
	if err != nil {
		return err

	}

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No student found with id %d", id)
	}

	return nil
}
