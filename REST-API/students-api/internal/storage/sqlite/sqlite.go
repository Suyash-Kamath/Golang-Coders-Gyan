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



func New(cfg *config.Config) (*Sqlite,error){
	// Open database , means connection
	// database ka driver pass and storage path
	db,err:=sql.Open("sqlite3",cfg.StoragePath)
	if err!=nil{
		return nil,err
	}

	_,err =db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)
	// hey _,err:= me why did you remove the colon

	if err!=nil{
		return nil,err
	}

	return &Sqlite{
		Db: db,
	},nil

}


func (s *Sqlite) CreateStudent(name string,email string,age int) (int64,error){

	stmt,err:=s.Db.Prepare("INSERT INTO students(name,email,age) VALUES (?,?,?)")
	// we dont pass the data directly , first we prepare the query and give the placeholder and then we bind the data , which prevenets the SQL injection attack

	if err!=nil{
		return 0,err
	}
	// stmt ko bhi close karna hota hai 
	defer stmt.Close()

	// as we have prepared the query so execute bhi karna hai 

	result,err:=stmt.Exec(name,email,age) // binds the stmt and replaces the placeholder

	if err!=nil{
		return 0,err
	}

	// result ke andhar , seeing the documentation , we get LastInsertID() , RowsAffected() function

	lastId,err:=result.LastInsertId()
	if err!=nil{
		return 0,err
		// int64 ka empty value is 0
	}
	return lastId,nil

}

func (s *Sqlite) GetStudentById(id int64) (types.Student,error){

	// stmt,err:=s.Db.Prepare(`SELECT * FROM students WHERE id = ? LIMIT 1`)
	stmt,err:=s.Db.Prepare(`SELECT id,name,email,age FROM students WHERE id = ? LIMIT 1`)
	if err!=nil{
		return types.Student{

		},err
	}

	defer stmt.Close()

	// data jo database se , aar aha hai , struct ke andhar serialize karke daaldo

	var student types.Student
	// Scan method database ka data , struct me dalega , kya kya scan karna hai wo batana hai , and order wise hona chahiye .
	err=stmt.QueryRow(id).Scan(&student.Id,&student.Name,&student.Email,&student.Age)
	if err!=nil{
		// if no user found
		if err == sql.ErrNoRows{
			return types.Student{} , fmt.Errorf("No student found with id %s",fmt.Sprint(id))
		}
		return types.Student{},fmt.Errorf("Query error %w",err)
	}


	return student,nil
}